package datastore

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/entity"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"log"
	"testing"
)

func TestPost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          entity.Employee
		expectedOutput entity.Employee
		err            error
	}{
		{"valid input", entity.Employee{"123", "jason", "Bangalore", "CSE"}, entity.Employee{"123", "jason", "Bangalore", "CSE"}, nil},
		{"Db error", entity.Employee{"123", "jason", "Bangalore", "CSE"}, entity.Employee{}, errors.New("error in query")},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.DataStore.DB().DB = db
	for i, tc := range testcases {
		mock.ExpectExec("Insert into emp values").
			WithArgs(tc.input.ID, tc.input.Name, tc.input.City, tc.input.Majors).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tc.err)
		actualOutput, _ := New().Post(ctx, tc.input)

		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestPut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             string
		input          entity.Employee
		expectedOutput entity.Employee
		err            error
	}{
		{"valid input", "123", entity.Employee{"123", "roy", "Kochi", "MBA"}, entity.Employee{"123", "roy", "Kochi", "MBA"}, nil},
		{"db error", "123", entity.Employee{"123", "roy", "Kochi", "MBA"}, entity.Employee{}, errors.New("error in db")},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.DataStore.DB().DB = db

	for i, tc := range testcases {
		mock.ExpectExec("Update").WithArgs(tc.input.ID, tc.input.Name, tc.input.City, tc.input.Majors, tc.id).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tc.err)
		actualOutput, _ := New().Put(ctx, tc.id, tc.input)

		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestDelete(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput int
		err            error
	}{
		{"id in db", "123", 204, nil},
		{"db error", "124", 400, errors.New("error in query")},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.DataStore.DB().DB = db

	for i, tc := range testcases {
		mock.ExpectExec("Delete").WithArgs(tc.input).WillReturnResult(sqlmock.
			NewResult(1, 1)).WillReturnError(tc.err)
		actualOutput, _ := New().Delete(ctx, tc.input)

		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}
