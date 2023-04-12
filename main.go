package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// region inicialização #############################
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

func configServer() { // TODO - Bloquear acesso no endpoint por outros metodos
	router := mux.NewRouter()
	router.HandleFunc("/heroi", novoHeroi).Methods("PUT")
	router.HandleFunc("/heroi", listaHerois).Methods("GET")
	router.HandleFunc("/heroi/{id}", listaHeroi).Methods("GET")
	router.HandleFunc("/heroi/{id}", excluiHeroi).Methods("DELETE")
	router.HandleFunc("/heroi/{id}", mudaHeroi).Methods("POST")
	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
}

// endregion inicialização

// region utilitários #############################
func processaBody(body io.ReadCloser) []byte {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return b
}

func checaPorId(w http.ResponseWriter, r *http.Request) string {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Indique o Id do Herói que deseja buscar!")
		return ""
	}
	return id
}

// endregion utilitários

// region CRUD Heroi #############################
func listaHerois(w http.ResponseWriter, r *http.Request) { // FIX - Gambiarra
	if len(r.URL.Query()) == 0 {
		hers, _ := buscaTodosHerois()
		js_hers, err := json.Marshal(hers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Erro ao procurar Herois!")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(js_hers))
	} else {
		hers, err := buscaHerois(r.URL.Query())
		if hers == nil || err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Nenhum Herói encontrado!")
			return
		}
		js_hers, err := json.Marshal(hers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Erro ao procurar Herois!")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(js_hers))
	}
}

func buscaTodosHerois() ([]Heroi, error) {
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
	return herois, nil
}

func buscaHerois(v url.Values) ([]Heroi, error) {
	var herois []Heroi
	var aux []Heroi
	// Busca por nome
	if v["nome"] != nil {
		rows, err := db.Query("SELECT * FROM heroi WHERE nome LIKE ?", "%"+v["nome"][0]+"%")
		if err != nil {
			return nil, fmt.Errorf("buscaHerois %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			var her Heroi
			if err := rows.Scan(&her.Id, &her.Nome, &her.Classe, &her.Ranking); err != nil {
				return nil, fmt.Errorf("todosHerois %v", err)
			}
			aux = append(aux, her)
		}
	}
	// Busca por outros campos
	keys := make([]string, len(v))
	i := 0
	for k := range v {
		keys[i] = k
		i++
	}
	for k := range keys {
		if keys[k] != "nome" {
			qy := fmt.Sprintf("SELECT * FROM heroi WHERE %s= ?", keys[k])
			rows, err := db.Query(qy, v[keys[k]][0])
			if err != nil {
				return nil, fmt.Errorf("buscaHerois %v", err)
			}
			defer rows.Close()
			for rows.Next() {
				var her Heroi
				if err := rows.Scan(&her.Id, &her.Nome, &her.Classe, &her.Ranking); err != nil {
					return nil, fmt.Errorf("todosHerois %v", err)
				}
				// adiciona heroi no array auxiliar sem duplicar
				found := false
				for _, h := range aux {
					if h == her {
						found = true
					}
				}
				if found == false {
					aux = append(aux, her)
				}
			}
		}
	}
	// checa se da match em todos os parametros de busca e adiciona no array de retorno se sim
	for _, h := range aux { // FIX - Gambiarra
		if v["nome"] != nil {
			if !strings.Contains(h.Nome, v["nome"][0]) {
				continue
			}
		}
		if v["classe"] != nil {
			if v["classe"][0] != h.Classe {
				continue
			}
		}
		if v["ranking"] != nil {
			r, _ := strconv.Atoi(v["ranking"][0])
			if r != h.Ranking {
				continue
			}
		}
		herois = append(herois, h)
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
	hers, err := buscaTodosHerois()
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

func mudaHeroi(w http.ResponseWriter, r *http.Request) {
	id := checaPorId(w, r)
	if id == "-1" {
		return
	}
	data := processaBody(r.Body)
	var h Heroi
	json.Unmarshal(data, &h)
	_id, _ := strconv.Atoi(id)
	if h.Classe != "" {
		_, err := db.Exec("UPDATE heroi SET classe=? WHERE id=?", h.Classe, _id)
		if err != nil {
			panic(err.Error())
		}
	}
	if h.Ranking != 0 {
		_, err := db.Exec("UPDATE heroi SET ranking=? WHERE id=?", h.Ranking, _id)
		if err != nil {
			panic(err.Error())
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func excluiHeroi(w http.ResponseWriter, r *http.Request) {
	id := checaPorId(w, r)
	if id == "-1" {
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

// endregion CRUD Heroi
