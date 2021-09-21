package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb_Error(t *testing.T) {
	err := DB{}

	assert.Nil(t, err.Err)
}
