package store

import (
	"context"
	"reflect"
	"testing"

	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func TestGetSetDelete(t *testing.T) {
	k := gofr.New()
	c := gofr.NewContext(nil, nil, k)
	c.Context = context.Background()

	//initializing the seeder
	seeder := datastore.NewSeeder(&k.DataStore, "../db")
	seeder.RefreshRedis(t, "store")

	testSet(t, c)
	testGet(t, c)
	testDelete(t, c)
	testSetWithError(t, k, c)
}

func testSetWithError(t *testing.T, k *gofr.Gofr, c *gofr.Context) {
	k.Redis.Close()

	expected := "redis: client is closed"
	got := Model{}.Set(c, "key", "value", 0)

	if !reflect.DeepEqual(expected, got.Error()) {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, got)
	}
}

func testSet(t *testing.T, c *gofr.Context) {
	tests := []struct {
		key         string
		value       string
		expectedErr error
	}{
		{
			key:         "someKey123",
			value:       "someValue123",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		err := Model{}.Set(c, test.key, test.value, 0)

		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, err)
		}
	}
}

func testGet(t *testing.T, c *gofr.Context) {
	tests := []struct {
		key      string
		expected string
		err      error
	}{
		{
			key:      "someKey123",
			expected: "someValue123",
			err:      nil,
		},
		{
			key:      "someKey",
			expected: "",
			err:      errors.DB{},
		},
	}

	for i, test := range tests {
		got, err := Model{}.Get(c, test.key)

		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expected, got)
		}

		if test.err == nil {
			if err != nil {
				t.Errorf("Testcase: %v FAILED", i)
			}
		} else {
			if _, ok := err.(errors.DB); ok == false {
				t.Errorf("Testcase: %v FAILED", i)
			}
		}
	}
}

func testDelete(t *testing.T, c *gofr.Context) {
	tests := []struct {
		key         string
		expectedErr error
	}{
		{
			key:         "someKey123",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		err := Model{}.Delete(c, test.key)

		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, err)
		}
	}
}

func TestNew(t *testing.T) {
	_ = New()
}
