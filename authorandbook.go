package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	driver            = "mysql"
	datasource        = "root:mukheshM1@25-03@/testingDB"
	createDatabase    = "CREATE DATABASE testingDB"
	createAuthorTable = "CREATE TABLE IF NOT EXISTS Author(authorId int NOT NULL AUTO_INCREMENT,Firstname varchar(50),Lastname varchar(50),Dob varchar(50),Penname varchar(50),PRIMARY KEY(authorId))"
	createBookTable   = "CREATE TABLE IF NOT EXISTS Book(bookId int NOT NULL AUTO_INCREMENT,Title varchar(50),authorId int,Publication varchar(50),PublishedDate varchar(50),PRIMARY KEY(bookId),CONSTRAINT FK FOREIGN KEY(authorId) REFERENCES Author(authorId))"
	GetAllBooks       = "SELECT * FROM Book"
	GetOneBook        = "SELECT * FROM Book WHERE bookId = ?"
)

type book struct {
	BookId        int    `json:"bookId"`
	Title         string `json:"title"`
	Author        author `json:"author"`
	Publication   string `json:"publication"`
	PublishedDate string `json:"publishedDate"`
}
type author struct {
	AuthorId  int    `json:"authorId"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Dob       string `json:"dob"`
	Penname   string `json:"penName"`
}

// Createtable establishes database connection and to create tables
func Createtable() (*sql.DB, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("USE testingDB")
	if err != nil {
		log.Print(err)
	}
	_, err = db.Exec(createAuthorTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(createBookTable)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

//PublishedDate used to check whether the date is valid or not
func PublishedDate(PublishedDate string) bool {

	date := strings.Split(PublishedDate, "/")
	year, _ := strconv.Atoi(date[0])
	month, _ := strconv.Atoi(date[1])
	day, _ := strconv.Atoi(date[2])
	switch {
	case year <= 1880 || year > 2022:
		return false
	case month <= 0 || month > 12:
		return false
	case day <= 0 || day > 31:
		return false
	}
	return true
}

func Publications(publications string) bool {

	switch strings.ToLower(publications) {
	case "penguin":
		return true
	case "scholastic":
		return true
	case "arihant":
		return true
	default:
		return false
	}
}

func checkDob(DOB string) bool {
	if DOB == "" {
		return false
	}
	dob := strings.Split(DOB, "/")
	year, _ := strconv.Atoi(dob[0])
	month, _ := strconv.Atoi(dob[1])
	day, _ := strconv.Atoi(dob[2])
	switch {
	case year > 2022:
		return false
	case month <= 0 || month > 12:
		return false
	case day <= 0 || day > 31:
		return false
	}
	return true
}

func GetbyId(w http.ResponseWriter, req *http.Request) {
	db, _ := Createtable()
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if id <= 0 {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bookrow := db.QueryRow(GetOneBook, id)
	b := book{}
	err = bookrow.Scan(&b.BookId, &b.Title, &b.Author.AuthorId, &b.Publication, &b.PublishedDate)
	if err != nil {
		log.Print(err)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	authorrow := db.QueryRow("SELECT * FROM Author WHERE authorId=?", b.Author.AuthorId)
	err = authorrow.Scan(&b.Author.AuthorId, &b.Author.Firstname, &b.Author.Lastname, &b.Author.Dob, &b.Author.Penname)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func GetAll(w http.ResponseWriter, req *http.Request) {
	db, err := Createtable()
	if err != nil {
		log.Fatal(err)
	}
	title := req.URL.Query().Get("title")
	includeAuthor := req.URL.Query().Get("includeAuthor")
	var rows *sql.Rows
	if title == "" {
		rows, err = db.Query(GetAllBooks)
	} else {
		rows, err = db.Query("select * from Books where title=?;", title)
	}
	if err != nil {
		log.Fatal(err)
	}
	books := []book{}
	for rows.Next() {
		b := book{}
		err := rows.Scan(&b.BookId, &b.Title, &b.Author, &b.Publication, &b.PublishedDate)
		if err != nil {
			log.Fatal(err)
		}
		if includeAuthor == "true" {
			row := db.QueryRow("select * from Author where authorId =?", b.Author.AuthorId)
			row.Scan(&b.Author.AuthorId, &b.Author.Firstname, &b.Author.Lastname, &b.Author.Dob, &b.Author.Penname)
		}
		books = append(books, b)
		mb, err := json.Marshal(books)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(mb)
	}

}

func PostBook(w http.ResponseWriter, req *http.Request) {
	db, err := Createtable()
	if err != nil {
		log.Print(err)
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
	}

	var b book
	json.Unmarshal(body, &b)
	if b.BookId <= 0 || b.Title == "" || b.Author.AuthorId <= 0 {
		fmt.Println("invalid constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !PublishedDate(b.PublishedDate) || !Publications(b.Publication) {
		fmt.Println("invalid constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = db.Exec("INSERT INTO Book(bookId,Title,authorId,Publication,PublishedDate) VALUES(?,?,?,?,?)", &b.BookId, &b.Title, &b.Author.AuthorId, &b.Publication, &b.PublishedDate)
	if err != nil {
		log.Print(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}
func PostAuthor(w http.ResponseWriter, req *http.Request) {
	db, _ := Createtable()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	var a author
	json.Unmarshal(body, &a)
	if a.AuthorId <= 0 || a.Firstname == "" || a.Lastname == "" || a.Penname == "" {
		fmt.Println("invalid constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !checkDob(a.Dob) {
		fmt.Println("invalid constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO Author(authorId,Firstname,Lastname,Dob,Penname) VALUES(?,?,?,?,?)", &a.AuthorId, &a.Firstname, &a.Lastname, &a.Dob, &a.Penname)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func PutBook(w http.ResponseWriter, req *http.Request) {
	db, _ := Createtable()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	var b book
	json.Unmarshal(body, &b)
	if b.BookId <= 0 || b.Title == "" || b.Author.AuthorId <= 0 {
		fmt.Println("invalid constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !checkDob(b.Author.Dob) || !PublishedDate(b.PublishedDate) || !Publications(b.Publication) {
		fmt.Println("invalid constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row, err := db.Query("SELECT * FROM Book WHERE bookId=?", b.BookId)
	if err != nil {
		log.Fatal(err)
	}
	if !row.Next() {
		log.Print("no id found")
		w.WriteHeader(http.StatusBadRequest)
	}
	row, err = db.Query("UPDATE Book SET Title=?,authorId=?,Publication=?,PublishedDate=? WHERE bookId=?", &b.Title, &b.Author.AuthorId, &b.Publication, &b.PublishedDate, &b.BookId)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(body)

}
func PutAuthor(w http.ResponseWriter, req *http.Request) {
	db, _ := Createtable()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	var a author
	json.Unmarshal(body, &a)
	if a.AuthorId <= 0 || a.Firstname == "" || a.Lastname == "" || a.Penname == "" {
		fmt.Println("invaid constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !checkDob(a.Dob) {
		fmt.Println("invalid constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row, err := db.Query("SELECT * FROM Author WHERE authorId =?", a.AuthorId)
	if err != nil {
		log.Fatal(err)
	}
	if !row.Next() {
		log.Print("author id not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row, err = db.Query("UPDATE Author SET Firstname=?,Lastname=?,Dob=?,Penname=? WHERE authorId=?", &a.Firstname, &a.Lastname, &a.Dob, &a.Penname, &a.AuthorId)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(body)

}

func DeleteBook(w http.ResponseWriter, req *http.Request) {
	db, err := Createtable()
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])
	if id <= 0 {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM Book WHERE bookId =?", id)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)
}
func DeleteAuthor(w http.ResponseWriter, req *http.Request) {
	db, err := Createtable()
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])
	if id <= 0 {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM Author WHERE authorId =?", id)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)

}
