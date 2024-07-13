package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/klaidliadon/trace/proto"
	"github.com/klaidliadon/trace/rpc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

func main() {
	logger := httplog.NewLogger("trace-id-example", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		MessageFieldName: "message",
		QuietDownPeriod:  10 * time.Second,
		SourceFieldName:  "source",
		Trace: &httplog.TraceOptions{
			LogFieldTrace: "logging.googleapis.com/trace",
			LogFieldSpan:  "logging.googleapis.com/spanId",
		},
	})

	s := Server{
		logger: logger,
	}

	httpClient := &http.Client{Transport: httplog.NewTransport("", nil)}

	portTer, err := s.Run(proto.NewTertiaryServer(rpc.Tertiary{}))
	if err != nil {
		logger.Error("failed to run tertiary server", slog.Any("error", err))
		return
	}

	clientTer := proto.NewTertiaryClient(fmt.Sprintf("http://localhost:%d", portTer), httpClient)

	portSec, err := s.Run(proto.NewSecondaryServer(rpc.Secondary{Tertiary: clientTer}))
	if err != nil {
		logger.Error("failed to run secondary server", slog.Any("error", err))
		return
	}

	clientSec := proto.NewSecondaryClient(fmt.Sprintf("http://localhost:%d", portSec), httpClient)

	portMain, err := s.Run(proto.NewMainServer(rpc.Main{Secondary: clientSec}))
	if err != nil {
		logger.Error("failed to run main server", slog.Any("error", err))
		return
	}

	clientMain := proto.NewMainClient(fmt.Sprintf("http://localhost:%d", portMain), httpClient)

	if _, err := clientTer.Do(context.Background()); err != nil {
		logger.Error("Secondary client", slog.Any("error", err))
	}
	if _, err := clientSec.Do(context.Background()); err != nil {
		logger.Error("Secondary client", slog.Any("error", err))
	}
	if _, err := clientMain.Do(context.Background()); err != nil {
		logger.Error("Main client", slog.Any("error", err))
	}
}

type Server struct {
	logger *httplog.Logger
}

func (s Server) Run(svc http.Handler) (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("listen: %w", err)
	}

	addr := l.Addr().String()

	logger := s.logger.With(slog.String("service", reflect.TypeOf(svc).String()))
	logger.Info("listening", slog.String("address", addr))

	go func() {
		r := chi.NewRouter()
		r.Use(httplog.RequestLogger(s.logger))
		r.Handle("/*", svc)
		if err := http.Serve(l, r); err != nil {
			if opError := (&net.OpError{}); errors.As(err, &opError) {
				logger.Error("server", slog.Any("error", opError.Err))
			}
		}
	}()

	return l.Addr().(*net.TCPAddr).Port, nil
}
