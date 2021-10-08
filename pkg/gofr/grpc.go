package gofr

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"developer.zopsmart.com/go/gofr/pkg/log"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"

	"google.golang.org/grpc"
)

type GRPC struct {
	server *grpc.Server
	Port   int
}

func (g *GRPC) Server() *grpc.Server {
	return g.server
}

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			LoggingInterceptor(log.NewLogger()),
		)))
}

func (g *GRPC) Start(logger log.Logger) {
	addr := ":" + strconv.Itoa(g.Port)

	logger.Infof("starting grpc server at %s", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorf("error in starting grpc server at %s: %s", addr, err)
		return
	}

	if err := g.server.Serve(listener); err != nil {
		logger.Errorf("error in starting grpc server at %s: %s", addr, err)
		return
	}
}

type RPCLog struct {
	ID           string `json:"id"`
	StartTime    string `json:"startTime"`
	ResponseTime int64  `json:"responseTime"`
	Method       string `json:"method"`
}

func (l RPCLog) String() string {
	line, _ := json.Marshal(l)
	return string(line)
}

func LoggingInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		ctx, span := trace.StartSpan(ctx, info.FullMethod)
		start := time.Now()

		defer func() {
			l := RPCLog{
				ID:           trace.FromContext(ctx).SpanContext().TraceID.String(),
				StartTime:    start.Format("2006-01-02T15:04:05.999999999-07:00"),
				ResponseTime: time.Since(start).Microseconds(),
				Method:       info.FullMethod,
			}

			if logger != nil {
				logger.Infof("%s", l)
			}

			span.End()
		}()

		return handler(ctx, req)
	}
}
