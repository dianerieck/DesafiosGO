package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Cotacao struct {
	ID  int
	Bid string
}

func NovaCotacao(bid string) *Cotacao {
	return &Cotacao{
		Bid: bid,
	}
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "./cotacoes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := r.Context()
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")

	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Buscando cotação...")
		req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		var cotacoes map[string]Cotacao
		err = json.NewDecoder(res.Body).Decode(&cotacoes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		novaCotacao := NovaCotacao(cotacoes["USDBRL"].Bid)

		// salva no banco
		err = salvarCotacaoNoBanco(db, novaCotacao)
		if err != nil {
			panic(err)
		}

		p, err := buscaCotacao(db)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"bid": p.Bid})

	case <-ctx.Done():
		log.Println("Request cancelada pelo cliente")
	}
}

func salvarCotacaoNoBanco(db *sql.DB, cotacao *Cotacao) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT)`)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := db.Prepare("INSERT INTO cotacoes(bid) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cotacao.Bid)
	if err != nil {
		return err
	}
	return nil
}

func buscaCotacao(db *sql.DB) (*Cotacao, error) {

	stmt, err := db.Prepare(`SELECT id, bid FROM cotacoes order by id DESC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var b Cotacao

	err = stmt.QueryRow().Scan(&b.ID, &b.Bid)
	if err != nil {
		return nil, err
	}
	return &b, nil

}
