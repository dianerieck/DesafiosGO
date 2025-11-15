package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	valor := result["bid"]
	if valor == "" {
		panic("Valor da cotação não encontrado")
	}
	err = os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", valor)), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Arquivo cotacao.txt salvo com sucesso!")
}
