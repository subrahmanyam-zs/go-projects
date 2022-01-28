package gofr

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/log"

	"google.golang.org/grpc"
)

func TestRPCLog_String(t *testing.T) {
	l := RPCLog{
		ID:           "123",
		StartTime:    "2020-01-01T12:12:12",
		ResponseTime: 100,
		Method:       http.MethodGet,
	}

	expected := `{"id":"123","startTime":"2020-01-01T12:12:12","responseTime":100,"method":"GET"}`
	got := l.String()

	if got != expected {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, got)
	}
}

func TestGRPC_Server(t *testing.T) {
	tcs := []struct {
		input *grpc.Server
	}{
		{nil},
		{new(grpc.Server)},
	}

	for _, tc := range tcs {
		g := new(GRPC)
		g.server = tc.input

		if g.Server() != tc.input {
			t.Errorf("FAILED, Expected: %v, Got: %v", tc.input, g.Server())
		}
	}
}

func TestNewGRPCServer(t *testing.T) {
	g := NewGRPCServer()
	if g == nil {
		t.Errorf("FAILED, Expected: a non nil value, Got: %v", g)
	}
}

func TestGRPC_Start(t *testing.T) {
	type fields struct {
		server *grpc.Server
		Port   int
	}

	type args struct {
		logger log.Logger
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		expectedLog string
	}{
		{
			name:        "net.Listen() error",
			fields:      fields{server: nil, Port: 99999},
			expectedLog: "error in starting grpc server",
		},
		{
			name:        "server.Serve() error",
			fields:      fields{server: new(grpc.Server), Port: 10000},
			expectedLog: "error in starting grpc server",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			tt.args.logger = log.NewMockLogger(b)

			g := &GRPC{
				server: tt.fields.server,
				Port:   tt.fields.Port,
			}

			g.Start(tt.args.logger)

			if !strings.Contains(b.String(), "error in starting grpc server") {
				t.Errorf("FAILED, Expected: `%v` in logs", "error in starting grpc server")
			}
		})
	}
}
