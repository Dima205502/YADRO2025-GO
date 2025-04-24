package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func uploadFilesHandler(w http.ResponseWriter, r *http.Request, cfg *Config) {
	file, header, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("uploadFilesHandler", "r.FormFile", err.Error())
		return
	}
	defer file.Close()

	filePath := filepath.Join(cfg.StoragePath, header.Filename)

	if _, err := os.Stat(filePath); err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("uploadFilesHandler", "os.Create", err.Error())
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("uploadFilesHandler", "io.Copy", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(header.Filename))
	if err != nil {
		slog.Error("uploadFilesHandler", "w.Write", err.Error())
	}
}

func replaceFilesHandler(w http.ResponseWriter, r *http.Request, cfg *Config) {
	filename := r.PathValue("filename")

	filePath := filepath.Join(cfg.StoragePath, filename)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("replaceFilesHandler", "os.Stat", err.Error())
		}
		return
	}

	file, _, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("replaceFilesHandler", "r.FormFile", err.Error())
		return
	}
	defer file.Close()

	dstFile, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("replaceFilesHandler", "os.Create", err.Error())
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("replaceFilesHandler", "io.Copy", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func listFilesHandler(w http.ResponseWriter, r *http.Request, cfg *Config) {
	entries, err := os.ReadDir(cfg.StoragePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("listFilesHandler", "os.ReadDir", err.Error())
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			_, err := w.Write([]byte(entry.Name() + "\n"))
			if err != nil {
				slog.Error("listFilesHandler", "w.Write", err.Error())
			}
		}
	}
}

func downloadFilesHandler(w http.ResponseWriter, r *http.Request, cfg *Config) {
	filename := r.PathValue("filename")

	filePath := filepath.Join(cfg.StoragePath, filename)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("downloadFilesHandler", "os.Stat", err.Error())
		}
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("downloadFilesHandler", "os.Open", err.Error())
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("downloadFilesHandler", "io.Copy", err.Error())
		return
	}
}

func deleteFilesHandler(w http.ResponseWriter, r *http.Request, cfg *Config) {
	filename := r.PathValue("filename")
	filePath := filepath.Join(cfg.StoragePath, filename)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("deleteFilesHandler", "os.Stat", err.Error())
		}
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("deleteFilesHandler", "os.Remove", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
