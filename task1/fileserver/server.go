package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func main() {
	cfg, err := InitConfig()
	if err != nil {
		slog.Error("main", "InitConfig", err.Error())
	}

	if _, err := os.Stat(cfg.StoragePath); err != nil {
		err := os.Mkdir(cfg.StoragePath, 0777)
		if err != nil {
			slog.Error("main", "os.Stat", err.Error())
			return
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /files", func(w http.ResponseWriter, r *http.Request) {
		uploadFilesHandler(w, r, cfg)
	})

	mux.HandleFunc("PUT /files/{filename}", func(w http.ResponseWriter, r *http.Request) {
		replaceFilesHandler(w, r, cfg)
	})

	mux.HandleFunc("GET /files", func(w http.ResponseWriter, r *http.Request) {
		listFilesHandler(w, r, cfg)
	})

	mux.HandleFunc("GET /files/{filename}", func(w http.ResponseWriter, r *http.Request) {
		downloadFilesHandler(w, r, cfg)
	})

	mux.HandleFunc("DELETE /files/{filename}", func(w http.ResponseWriter, r *http.Request) {
		deleteFilesHandler(w, r, cfg)
	})

	adr := ":" + strconv.Itoa(cfg.Port)
	if err := http.ListenAndServe(adr, mux); err != nil {
		slog.Error("Сервер не запустился", "message", err.Error())
	}
}
