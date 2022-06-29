package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestGetbyId(t *testing.T) {
	testcase := []struct {
		desc          string
		inputurl      string
		expstatuscode int
		expoutput     book
	}{
		{"valid id", "1", http.StatusOK, book{1, "it",
			author{1, "chetan", "bhagat", "1980/04/23", "rolex"},
			"arihant", "2020/10/23"}},
		{"invalid id", "-1", http.StatusBadRequest, book{}},
	}
	for _, tc := range testcase {
		req := httptest.NewRequest("GET", "/books/{id}"+tc.inputurl, nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.inputurl})
		w := httptest.NewRecorder()
		GetbyId(w, req)
		res := w.Result()
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var b book
		json.Unmarshal(body, &b)

		assert.Equal(t, tc.expoutput, b)

	}

}

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc       string
		url        string
		expetedout []book
	}{
		{"getting data of all books", "/books",
			[]book{{1, "it", author{1, "chetan", "bhagat", "1980/04/23", "rolex"},
				"arihant", "2020/10/23"},
				book{2, "god", author{2, "mukhesh", "mekala", "2000/10/24", "zero"},
					"arihant", "2015/07/03"}},
		},
	}
	for _, tc := range testcases {
		req := httptest.NewRequest("GET", "/books", nil)
		w := httptest.NewRecorder()
		GetAll(w, req)
		res := w.Result()
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var b book
		json.Unmarshal(body, &b)

		assert.Equal(t, tc.expetedout, b)

	}

}

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc          string
		inputurl      string
		input         book
		expstatuscode int
	}{
		{"details posted for id 1", "/books", book{1, "2states", author{1, "chetan", "bhagat", "1980/04/23", "max"}, "arihant", "2020/10/23"}, http.StatusCreated},
		{"details posted for id 2", "/books", book{2, "god", author{2, "mukhesh", "mekala", "2000/10/24", "zero"}, "arihant", "2015/07/03"}, http.StatusCreated},
		{"invalid input title", "/books", book{2, "", author{2, "mukhesh", "mekala", "2000/10/24", "zero"}, "arihant", "2015/07/03"}, http.StatusBadRequest},
	}
	for _, tc := range testcases {
		data, _ := json.Marshal(tc.input)
		req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		PostBook(w, req)
		res := w.Result()
		assert.Equal(t, tc.expstatuscode, res.StatusCode)
	}

}
func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc          string
		input         author
		expstatuscode int
	}{
		{"valid author", author{1, "chetan", "bhagat", "1980/04/23", "max"}, http.StatusCreated},
		{"valid author", author{2, "mukhesh", "mekala", "2000/10/24", "zero"}, http.StatusCreated},
		{"author firstname is empty", author{2, "", "kumar", "1906/06/16", "rolex"}, http.StatusBadRequest},
		{"author lastname is empty", author{2, "kumar", "", "1906/06/16", "rolex"}, http.StatusBadRequest},
		{"empty Dob", author{2, "kumar", "krishna", "", "rolex"}, http.StatusBadRequest},
		{"invalid author id", author{-2, "kumar", "krishna", "1906/06/16", "rolex"}, http.StatusBadRequest},
		{"penname not given", author{2, "kumar", "krishna", "1906/06/16", ""}, http.StatusBadRequest},
	}
	for _, tc := range testcases {
		data, _ := json.Marshal(tc.input)
		req := httptest.NewRequest("POST", "/author", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		PostAuthor(w, req)
		res := w.Result()
		assert.Equal(t, tc.expstatuscode, res.StatusCode)
	}
}

func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc          string
		inputurl      string
		input         book
		expstatuscode int
	}{
		{"updating details", "1", book{1, "it", author{1, "chetan", "bhagat", "1980/04/23", "max"}, "arihant", "2020/10/23"}, http.StatusAccepted},
		{"Book id not found", "10000", book{10000, "it", author{1, "chetan", "bhagat", "1980/04/23", "max"}, "arihant", "2020/10/23"}, http.StatusBadRequest},
	}
	for _, tc := range testcases {
		body, _ := json.Marshal(tc.input)
		req := httptest.NewRequest("PUT", "/books/{id}"+tc.inputurl, bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		PutBook(w, req)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, tc.expstatuscode, res.StatusCode)
	}
}
func TestPutAuthor(t *testing.T) {
	testcases := []struct {
		desc          string
		inputurl      string
		input         author
		expstatuscode int
	}{
		{"updating details for id1", "1", author{1, "chetan", "bhagat", "1980/04/23", "rolex"}, http.StatusAccepted},
		{"invalid id", "-1", author{-1, "chetan", "bhagat", "1980/04/23", "rolex"}, http.StatusBadRequest},
		{"firstname is empty", "1", author{1, "", "bhagat", "1980/04/23", "rolex"}, http.StatusBadRequest},
		{"lastname is empty", "1", author{1, "chetan", "", "1980/04/23", "rolex"}, http.StatusBadRequest},
		{"empty date od birth", "1", author{1, "chetan", "bhagat", "", "rolex"}, http.StatusBadRequest},
		{"Penname is empty", "1", author{1, "chetan", "bhagat", "1980/04/23", ""}, http.StatusBadRequest},
	}
	for _, tc := range testcases {
		body, _ := json.Marshal(tc.input)
		req := httptest.NewRequest("PUT", "/author/{id}"+tc.inputurl, bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		PutAuthor(w, req)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, tc.expstatuscode, res.StatusCode)
	}
}

func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc               string
		inputurl           string
		expectedstatuscode int
	}{
		{"valid id", "1", http.StatusNoContent},
		{"valid id", "-4", http.StatusBadRequest},
	}
	for _, tc := range testcases {
		req := httptest.NewRequest("DELETE", "/books/{id}"+tc.inputurl, nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.inputurl})
		w := httptest.NewRecorder()
		DeleteBook(w, req)
		res := w.Result()
		assert.Equal(t, tc.expectedstatuscode, res.StatusCode)
	}

}
func TestDeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc               string
		inputurl           string
		expectedstatuscode int
	}{
		{"valid id", "4", http.StatusNoContent},
		{"valid id", "-4", http.StatusBadRequest},
	}
	for _, tc := range testcases {
		req := httptest.NewRequest("DELETE", "/author/{id}"+tc.inputurl, nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.inputurl})
		w := httptest.NewRecorder()
		DeleteAuthor(w, req)
		res := w.Result()
		assert.Equal(t, tc.expectedstatuscode, res.StatusCode)
	}

}
