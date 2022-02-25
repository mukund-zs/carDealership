package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db *sql.DB

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:root@123@tcp(localhost:3306)/cardealership")
	if err != nil {
		log.Print("this is ", err)
		return
	}
	return db
}
func main() {
	r := mux.NewRouter()
	var err error
	r.HandleFunc("/", carHandler)
	r.HandleFunc("/{id}", runById).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
	err1 := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err1)
	}
}
