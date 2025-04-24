package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	searchpb "yadro.com/course/proto/search"
	"yadro.com/course/search/adapters/db"
	searchgrpc "yadro.com/course/search/adapters/grpc"
	"yadro.com/course/search/adapters/words"
	"yadro.com/course/search/config"
	"yadro.com/course/search/core"
)

func main() {

	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()
	cfg := config.MustLoad(configPath)

	// logger
	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting server")
	log.Debug("debug messages are enabled")

	// database adapter
	storage, err := db.New(log, cfg.DBAddress)
	if err != nil {
		log.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}

	// words adapter
	words, err := words.NewClient(cfg.WordsAddress, log)
	if err != nil {
		log.Error("failed create Words client", "error", err)
		os.Exit(1)
	}

	// service
	searcher := core.NewService(log, storage, words)

	if err := run(log, searcher, cfg); err != nil {
		log.Error("failed to serve", "erorr", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger, searcher *core.Service, cfg config.Config) error {

	// grpc server
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	searchpb.RegisterSearchServer(s, searchgrpc.NewServer(searcher))
	reflection.Register(s)

	// context for Ctrl-C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		s.GracefulStop()
	}()

	return s.Serve(listener)
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		panic("unknown log level: " + logLevel)
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
