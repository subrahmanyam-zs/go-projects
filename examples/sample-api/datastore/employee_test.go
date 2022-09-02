package datastore

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/entities"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestPost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          entities.Employee
		expectedOutput entities.Employee
	}{
		{"valid input", entities.Employee{"123", "jason", "Bangalore", "CSE"}, entities.Employee{"123", "jason", "Bangalore", "CSE"}},
	}

	db,mock,err :=
}
