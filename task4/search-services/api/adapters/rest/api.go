package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"yadro.com/course/api/core"
)

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}

type WordsResponse struct {
	Words []string `json:"words"`
	Total int      `json:"total"`
}

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("PingHandler", "pingers", pingers)

		replies := make(map[string]string)

		ctx := r.Context()
		for name, pinger := range pingers {
			if err := pinger.Ping(ctx); err != nil {
				replies[name] = "unavailable"

				log.Error("ping failed", "service", name, "error", err)
			} else {
				replies[name] = "ok"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(PingResponse{Replies: replies}); err != nil {
			log.Error("failed to encode JSON response", "error", err)
			return
		}
	}
}

func NewWordsHandler(log *slog.Logger, normalizer core.Normalizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("WordsHandlers", "phrase", r.URL.Query().Get("phrase"))

		phrase := r.URL.Query().Get("phrase")
		if phrase == "" {
			http.Error(w, "missing phrase query parameter", http.StatusBadRequest)
			return
		}

		words, err := normalizer.Norm(r.Context(), phrase)

		if err != nil {

			if errors.Is(err, core.ErrBadArguments) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			log.Error("failed to normalize phrase", "error", err, "phrase", phrase)
			http.Error(w, "failed to normalize phrase", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(WordsResponse{Words: words, Total: len(words)}); err != nil {
			log.Error("failed to encode JSON response", "error", err)
			return
		}
	}
}
