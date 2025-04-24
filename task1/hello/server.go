package main

import (
	"flag"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port int `yaml:"port" env:"HELLO_PORT"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong\n"))

	if err != nil {
		slog.Error("pingHandler", "w.Write", err.Error())
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("empty name\n"))

		if err != nil {
			slog.Error("helloHandler", "w.Write", err.Error())
		}
		return
	}

	_, err := w.Write([]byte("Hello, " + name + "!\n"))
	if err != nil {
		slog.Error("helloHandler", "w.Write", err.Error())
	}
}

func main() {
	configPath := flag.String("config", "", "Path to config.yaml file")
	flag.Parse()

	var cfg Config

	if *configPath != "" {

		err := cleanenv.ReadConfig(*configPath, &cfg)
		if err != nil {
			slog.Error("main", "cleanenv.ReadConfig", err.Error())
		}
	} else {

		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			slog.Error("main", "cleanenv.ReadEnv", err.Error())
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("GET /hello", helloHandler)

	adr := ":" + strconv.Itoa(cfg.Port)
	if err := http.ListenAndServe(adr, mux); err != nil {
		slog.Error("main", "http.ListenAndServe", err.Error())
	}
}
