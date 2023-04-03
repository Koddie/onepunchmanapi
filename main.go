package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type Heroi struct {
	Id      int    `json:"id"`
	Nome    string `json:"nome"`
	Classe  string `json:"classe"`
	Ranking int    `json:"ranking"`
}

var db *sql.DB

func main() {
	configDatabase()
	configServer()
}

func configDatabase() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "onepunchman",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Conectado ao Banco de Dados!")
}

func configServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/heroi/novo", novoHeroi)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func novoHeroi(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	reader := strings.NewReader(string(body))
	var h Heroi
	json.Unmarshal(body, &h)
	if err := json.NewDecoder(reader).Decode(&h); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Erro %s", err)
		return
	}
	// Erro se tenta cadastrar em classe não existente
	if h.Classe != "A" && h.Classe != "B" && h.Classe != "C" && h.Classe != "S" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "A classe do herói deve ser A, B, C ou S!")
		return
	}
	// Busca possiveis duplicacoes
	hers, err := listaTodosHerois()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "A busca por heróis já salvos falhou...")
		return
	}
	for _, her := range hers {
		if h.Nome == her.Nome {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Herói já cadastrado")
			return
		}
	}
	// Inserindo novo Herói
	result, err := db.Exec("INSERT INTO heroi (nome, classe, ranking) VALUES (?, ?, ?)", h.Nome, h.Classe, h.Ranking)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Erro ao salvar o novo herói...")
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Erro ao gerar o id do novo herói...")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{id: %d}", id)
}

func listaTodosHerois() ([]Heroi, error) {
	var herois []Heroi

	rows, err := db.Query("SELECT * FROM heroi")
	if err != nil {
		return nil, fmt.Errorf("todosHerois %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var her Heroi
		if err := rows.Scan(&her.Id, &her.Nome, &her.Classe, &her.Ranking); err != nil {
			return nil, fmt.Errorf("todosHerois %v", err)
		}
		herois = append(herois, her)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %v", err)
	}
	return herois, nil
}
