package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"yadro.com/course/frontend/config"
	"yadro.com/course/frontend/handler"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting server")
	log.Debug("debug messages are enabled")

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handler.HandlerRoot())

	mux.HandleFunc("GET /search", handler.HadlerSearch(http.DefaultClient, "http://"+cfg.Api_address, log))

	mux.HandleFunc("GET /login", handler.HandlerLogin())

	mux.HandleFunc("POST /login", handler.HandlerAuth(http.DefaultClient, "http://"+cfg.Api_address, log))

	mux.HandleFunc("GET /stats", handler.HandlerStats(http.DefaultClient, "http://"+cfg.Api_address, log))

	mux.HandleFunc("GET /status", handler.HandlerStatus(http.DefaultClient, "http://"+cfg.Api_address, log))

	mux.HandleFunc("GET /drop", handler.HandlerDrop(http.DefaultClient, "http://"+cfg.Api_address, log))

	mux.HandleFunc("GET /update", handler.HandlerUpdate(http.DefaultClient, "http://"+cfg.Api_address, log))

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
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
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level, AddSource: true})
	return slog.New(handler)
}
