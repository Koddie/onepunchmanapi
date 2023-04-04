package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.HandleFunc("/heroi/{id}", listaHeroi).Methods("GET")
	router.HandleFunc("/heroi/{id}", excluiHeroi).Methods("DELETE")
	router.HandleFunc("/heroi/novo", novoHeroi).Methods("POST")
	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
}

func processaBody(body io.ReadCloser) []byte {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return b
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

func listaHeroi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Indique o Id do Herói que deseja buscar!")
		return
	}
	result := db.QueryRow("SELECT * from heroi WHERE id = ?", id)
	var h Heroi
	if err := result.Scan(&h.Id, &h.Nome, &h.Classe, &h.Ranking); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{message: Herói não encontrado}")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao procurar o Herói!")
		return
	}
	str, err := json.Marshal(h)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao retornar o Herói!")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(str))
}

func novoHeroi(w http.ResponseWriter, r *http.Request) {
	data := processaBody(r.Body)
	var h Heroi
	json.Unmarshal(data, &h)
	// ERRO - se tenta cadastrar em classe não existente
	if h.Classe != "A" && h.Classe != "B" && h.Classe != "C" && h.Classe != "S" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "A classe do herói deve ser A, B, C ou S!")
		return
	}
	// ERRO - Se tenta salvar duplicado
	hers, err := listaTodosHerois()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao salvar o novo herói...")
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao gerar o id do novo herói...")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{id: %d}", id)
}

func excluiHeroi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Indique o Id do Herói que deseja buscar!")
		return
	}
	_, err := db.Exec("DELETE from heroi WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{message: Herói não encontrado}")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao excluir o Herói")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
