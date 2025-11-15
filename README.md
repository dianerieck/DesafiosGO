Projeto: Cotação do Dólar com Go

Este projeto implementa dois sistemas em Go que trabalham juntos para consultar a cotação do dólar (USD/BRL), persistir os dados em um banco SQLite e salvar o valor atual em um arquivo local.

- **server.go** → Servidor HTTP que consome a API externa de câmbio, persiste a cotação em SQLite e expõe o endpoint `/cotacao`.
- **client.go** → Cliente que consulta o servidor, aplica timeout, extrai o valor da cotação e salva em `cotacao.txt`.

---

Tecnologias utilizadas
- **Go** (>= 1.20)
- **SQLite** (via driver `modernc.org/sqlite`)
- **Contextos (`context`)** para controle de timeout
- **HTTP** para comunicação cliente-servidor

---

Como executar o projeto

1. Clonar o repositório

git clone <https://github.com/dianerieck/DesafiosGO>
cd <DesafiosGo/cotacao>

2. Instalar dependências
go get modernc.org/sqlite

3. Executar o servidor
go run server.go

O servidor iniciará na porta 8080 e disponibilizará o endpoint:
http://localhost:8080/cotacao

4. Executar o cliente
go run client.go

O cliente fará uma requisição ao servidor e:
Receberá o valor atual da cotação (campo bid).

Criará o arquivo cotacao.txt com o conteúdo:
Dólar: {valor}