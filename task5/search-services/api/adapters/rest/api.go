package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yadro.com/course/api/core"
)

type PingResponse struct {
	Replies map[string]string `json:"replies"`
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

func NewUpdateHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updater.Update(r.Context())

		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists && st.Message() == "already running" {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		if err != nil {
			http.Error(w, "NewUpdateHandler", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewUpdateStatsHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateStats, err := updater.Stats(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(updateStats); err != nil {
			log.Error("UpdateStatsHandlers", "error", err)
			return
		}
	}
}

func NewUpdateStatusHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		status, err := updater.Status(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(map[string]core.UpdateStatus{"status": status}); err != nil {
			log.Error("NewUpdateStatusHandler", "error", err)
		}
	}
}

func NewDropHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updater.Drop(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
