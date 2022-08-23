package main

import (
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/gofr/assert"
)

func TestIntegration(t *testing.T) {
	testCases := []struct {
		command  string
		expected string
	}{
		{"cmd read", "Welcome to Zopsmart!"},
		{"cmd write", "File written successfully!"},
		{"cmd list", "Readme.md configs handler main.go main_test.go test.txt"},
	}
	for _, tc := range testCases {
		assert.CMDOutputContains(t, main, tc.command, tc.expected)
	}
}
