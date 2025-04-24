package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yadro.com/course/api/core"
)

const defaultLimit = "10"

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

type Authenticator interface {
	Login(user, password string) (string, error)
}

func NewLoginHandler(log *slog.Logger, auth Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authInfo struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&authInfo); err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		token, err := auth.Login(authInfo.Name, authInfo.Password)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if _, err := w.Write([]byte(token)); err != nil {
			log.Error("NewLoginHandler", "w.Write", err)
		}

	}
}

func NewUpdateHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := updater.Update(r.Context()); err != nil {
			if code := status.Code(err); code == codes.AlreadyExists {
				w.WriteHeader(http.StatusAccepted)
				return
			}
			http.Error(w, "NewUpdateHandler:"+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewUpdateStatsHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := updater.Stats(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("StatsHandler", "result", stats)

		statsResponse := struct {
			WordsTotal    int `json:"words_total"`
			WordsUnique   int `json:"words_unique"`
			ComicsFetched int `json:"comics_fetched"`
			ComicsTotal   int `json:"comics_total"`
		}{
			WordsTotal:    stats.WordsTotal,
			WordsUnique:   stats.WordsUnique,
			ComicsFetched: stats.ComicsFetched,
			ComicsTotal:   stats.ComicsTotal,
		}

		if err := json.NewEncoder(w).Encode(statsResponse); err != nil {
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

type Comics struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type ComicsResponse struct {
	Comics []Comics `json:"comics"`
	Total  int      `json:"total"`
}

func NewSearchHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("SearchHandler start")

		query := r.URL.Query()
		limitStr := query.Get("limit")

		if limitStr == "" {
			limitStr = defaultLimit
		}

		limit, err := strconv.Atoi(limitStr)

		if err != nil || limit < 0 {
			http.Error(w, "Unexpected 'limit' parameter", http.StatusBadRequest)
			return
		}

		phrase := query.Get("phrase")
		if phrase == "" {
			http.Error(w, "Missing 'phrase' parameter.", http.StatusBadRequest)
			return
		}

		comics, err := searcher.DbSearch(r.Context(), limit, phrase)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Debug("Search", "limit", limit, "phrase", phrase)

		var comicsRespose ComicsResponse

		for _, x := range comics {
			comicsRespose.Comics = append(comicsRespose.Comics, Comics{ID: x.ID, URL: x.URL})
		}
		comicsRespose.Total = len(comicsRespose.Comics)

		log.Info("Search", "result", comicsRespose)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comicsRespose); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NewSearchIndexHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("IndexSearchHandler start")

		query := r.URL.Query()
		limitStr := query.Get("limit")

		if limitStr == "" {
			limitStr = defaultLimit
		}

		limit, err := strconv.Atoi(limitStr)

		if err != nil || limit < 0 {
			http.Error(w, "Unexpected 'limit' parameter", http.StatusBadRequest)
			return
		}

		phrase := query.Get("phrase")
		if phrase == "" {
			http.Error(w, "Missing 'phrase' parameter.", http.StatusBadRequest)
			return
		}

		comics, err := searcher.IndexSearch(r.Context(), limit, phrase)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Debug("IndexSearch", "limit", limit, "phrase", phrase)

		var comicsRespose ComicsResponse

		for _, x := range comics {
			comicsRespose.Comics = append(comicsRespose.Comics, Comics{ID: x.ID, URL: x.URL})
		}
		comicsRespose.Total = len(comicsRespose.Comics)

		log.Info("IndexSearch", "result", comicsRespose)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comicsRespose); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
