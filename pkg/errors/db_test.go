package errors

import (
	"errors"
	"testing"
)

func TestDb_Error(t *testing.T) {
	err := DB{}

	var expected error = nil
	if !errors.Is(expected, err.Err) {
		t.Errorf("FAILED Expected: %v, Got: %v", expected, err.Err)
	}
}
