package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"yadro.com/course/api/adapters/rest"
	"yadro.com/course/api/adapters/words"
	"yadro.com/course/api/config"
	"yadro.com/course/api/core"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.LoadConfig(configPath)
	fmt.Println(cfg)

	log := MakeLogger(cfg.LogLevel)

	log.Info("starting server")
	log.Debug("debug messages are enabled")

	wordsClient, err := words.NewClient(cfg.WordsAddress, log)
	if err != nil {
		log.Error("cannot init words adapter", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.Handle("GET /ping", rest.NewPingHandler(log, map[string]core.Pinger{"words": wordsClient}))
	mux.Handle("GET /api/words", rest.NewWordsHandler(log, wordsClient))

	server := http.Server{
		Addr:        cfg.Address,
		ReadTimeout: cfg.Timeout,
		Handler:     mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("erroneous shutdown", "error", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server closed unexpectedly", "error", err)
			return
		}
	}

	time.Sleep(cfg.Timeout)
}

func MakeLogger(logLevel string) *slog.Logger {
	var level slog.Level = slog.LevelDebug

	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		fmt.Printf("Unknown logging level %s\n", logLevel)
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: level}))
}
