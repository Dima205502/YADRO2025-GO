package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"yadro.com/course/frontend/model"
)

func HandlerRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HadlerSearch(client *http.Client, api_address string, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")
		limit := r.URL.Query().Get("limit")

		indexStr := r.URL.Query().Get("index")
		if indexStr == "" {
			indexStr = "0"
		}

		index, err := strconv.Atoi(indexStr)
		if err != nil {
			http.Error(w, "Не смогли достать index", http.StatusInternalServerError)
			return
		}

		log.Debug("HandlerSearch", "phrase", phrase, "limit", limit, "index", index)

		results, err := searchComics(client, api_address, phrase, limit)
		if err != nil {
			log.Error("HandlerSearch", "error", err)
			http.Error(w, "Не удалось найти картинки", http.StatusInternalServerError)
			return
		}

		data := model.TemplateData{
			Phrase:       phrase,
			Comics:       results.Comics,
			Total:        results.Total,
			CurrentIndex: index,
			DisplayTotal: len(results.Comics),
		}

		log.Debug("HandlerSearch", "data", data)

		tmpl := template.New("search.html").Funcs(template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"sub": func(a, b int) int { return a - b },
			"len": func(x interface{}) int {
				return reflect.ValueOf(x).Len()
			},
		})

		tmpl, err = tmpl.ParseFiles("templates/search/search.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlerLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/login/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlerAuth(client *http.Client, api_address string, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		data := model.AuthInfo{
			Name:     username,
			Password: password,
		}

		log.Debug("HandlerAuth", "username", username, "passwords", password)

		token, err := getToken(client, api_address, data)

		if err != nil {
			log.Error("HandlerAuth", "error", err)

			tmpl, err := template.ParseFiles("templates/auth/authorized_error.html")
			if err != nil {
				http.Error(w, "Не смогли открыть страничку 401", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		tmpl, err := template.ParseFiles("templates/login/success.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlerStats(client *http.Client, api_address string, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var stats model.StatsResponse

		resp, err := client.Get(api_address + "/api/db/stats")
		if err != nil {
			log.Error("HandlerStats", "error", err)
			http.Error(w, "Не удалось получить статистику", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			log.Error("HandlerStats", "error", err)
			http.Error(w, "Не удалось получить статистику", http.StatusInternalServerError)
			return
		}

		log.Debug("HandlerStats", "stats", stats)

		tmpl, err := template.ParseFiles("templates/stats/stats.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, stats)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlerStatus(client *http.Client, api_address string, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var status model.Status

		resp, err := client.Get(api_address + "/api/db/status")
		if err != nil {
			log.Error("HandlerStatus", "error", err)
			http.Error(w, "Не удалось получить статус", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			log.Error("HandlerStatus", "error", err)
			http.Error(w, "Не удалось получить статус", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/status/status.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlerDrop(client *http.Client, apiAddress string, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			tmpl, err := template.ParseFiles("templates/auth/unauthorized.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		req, err := http.NewRequest(http.MethodDelete, apiAddress+"/api/db", nil)
		if err != nil {
			http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Token "+cookie.Value)

		resp, err := client.Do(req)
		if err != nil {
			log.Error("HandlerDrop", "error", err)
			http.Error(w, "Ошибка при запросе к сервису удаления", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
			tmpl, err := template.ParseFiles("templates/auth/unauthorized.html")
			if err != nil {
				http.Error(w, "Не удалось открыть страницу недостаточно прав", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, nil)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/drop/drop.html")
		if err != nil {
			http.Error(w, "Не удалось открыть страницу успеха", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlerUpdate(client *http.Client, apiAddress string, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			tmpl, err := template.ParseFiles("templates/auth/unauthorized.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		req, err := http.NewRequest(http.MethodPost, apiAddress+"/api/db/update", nil)
		if err != nil {
			http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Token "+cookie.Value)

		resp, err := client.Do(req)
		if err != nil {
			log.Error("HandlerUpdate", "error", err)
			http.Error(w, "Ошибка при запросе к сервису обновления", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
			tmpl, err := template.ParseFiles("templates/auth/unauthorized.html")
			if err != nil {
				http.Error(w, "Не удалось открыть страницу недостаточно прав", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			http.Error(w, "Ошибка при обновлении", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/update/update.html")
		if err != nil {
			http.Error(w, "Не удалось открыть страницу успеха", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, struct {
			Code int
		}{Code: resp.StatusCode})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func searchComics(client *http.Client, api_address, phrase, limit string) (model.ComicsResponse, error) {
	resp, err := client.Get(api_address + "/api/search?" + "limit=" + limit + "&phrase=" + phrase)
	if err != nil {
		return model.ComicsResponse{}, err
	}

	var comicsResp model.ComicsResponse

	if err := json.NewDecoder(resp.Body).Decode(&comicsResp); err != nil {
		return model.ComicsResponse{}, err
	}

	return comicsResp, nil
}

func getToken(client *http.Client, api_address string, data model.AuthInfo) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, api_address+"/api/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("код ответа не 200")
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	token := string(bytes)

	return token, nil
}
