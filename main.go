package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books/{id}", GetbyId).Methods("GET")
	r.HandleFunc("/books", GetAll).Methods("GET")
	r.HandleFunc("/books", PostBook).Methods("POST")
	r.HandleFunc("/author", PostAuthor).Methods("POST")
	r.HandleFunc("/books/{id}", PutBook).Methods("PUT")
	r.HandleFunc("/author/{id}", PutAuthor).Methods("PUT")
	r.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")
	r.HandleFunc("/author/{id}", DeleteAuthor).Methods("DELETE")
	http.ListenAndServe(":8000", r)
}
