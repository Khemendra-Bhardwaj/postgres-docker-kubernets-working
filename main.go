package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	connStr := "user=postgres password=example dbname=postgres sslmode=disable host=postgres"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/get", GetHandler).Methods("GET")
	r.HandleFunc("/post", PostHandler).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, data FROM example_table")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		var id int
		var data string
		if err := rows.Scan(&id, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result += fmt.Sprintf("ID: %d, Data: %s\n", id, data)
	}
	fmt.Fprint(w, result)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	_, err := db.Exec("INSERT INTO example_table (data) VALUES ($1)", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Data inserted")
}
