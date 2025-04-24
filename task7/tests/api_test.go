package words_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const address = "http://localhost:28080"

var client = http.Client{
	Timeout: 5 * time.Minute,
}

func TestPreflight(t *testing.T) {
	require.Equal(t, true, true)
}

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}

func TestPing(t *testing.T) {
	resp, err := client.Get(address + "/api/ping")
	require.NoError(t, err, "cannot ping")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "wrong status")

	var reply PingResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&reply))
	require.Equal(t, "ok", reply.Replies["words"], "no words running")
	require.Equal(t, "ok", reply.Replies["update"], "no update running")
	require.Equal(t, "ok", reply.Replies["search"], "no search running")
}

type UpdateStats struct {
	WordsTotal    int `json:"words_total"`
	WordsUnique   int `json:"words_unique"`
	ComicsFetched int `json:"comics_fetched"`
	ComicsTotal   int `json:"comics_total"`
}

type UpdateStatus struct {
	Status string `json:"status"`
}

func prepare(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, address+"/api/db", nil)
	require.NoError(t, err, "cannot make request")
	resp, err := client.Do(req)
	require.NoError(t, err, "could not send clean up command")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	st := stats(t)
	require.Equal(t, 0, st.ComicsFetched)
	require.True(t, st.ComicsTotal > 3000, "there are more than 3000 comics in XKCD")
	require.Equal(t, 0, st.WordsTotal)
	require.Equal(t, 0, st.WordsUnique)

	require.Equal(t, "idle", status(t))
}

func TestEmptyDB(t *testing.T) {
	prepare(t)
}

func TestUpdate(t *testing.T) {
	prepare(t)
	var wg sync.WaitGroup
	wg.Add(3)
	var res1, res2 int
	var res3 string
	go func() {
		res1 = update(t)
		wg.Done()
	}()
	go func() {
		res2 = update(t)
		wg.Done()
	}()
	go func() {
		time.Sleep(1 * time.Second)
		res3 = status(t)
		wg.Done()
	}()
	wg.Wait()
	require.True(t,
		res1 == http.StatusOK && res2 == http.StatusAccepted ||
			res2 == http.StatusOK && res1 == http.StatusAccepted,
		"wrong statuses from concurrent updates, expect ok && accepted",
	)
	require.Equal(t, "running", res3, "need running status while update")
	st := stats(t)
	require.Equal(t, st.ComicsTotal, st.ComicsFetched)
	require.True(t, st.ComicsTotal > 3000, "there are more than 3000 comics in XKCD")
	require.True(t, 1000 < st.WordsTotal, "not enough total words in DB")
	require.True(t, 100 < st.WordsUnique, "not enough unique words in DB")
}

func update(t *testing.T) int {
	req, err := http.NewRequest(http.MethodPost, address+"/api/db/update", nil)
	require.NoError(t, err, "cannot make request")
	resp, err := client.Do(req)
	require.NoError(t, err, "could not send update command")
	defer resp.Body.Close()
	return resp.StatusCode
}

func status(t *testing.T) string {
	resp, err := client.Get(address + "/api/db/status")
	require.NoError(t, err, "could not get status")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var status UpdateStatus
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&status), "cannot decode")
	return status.Status
}

func stats(t *testing.T) UpdateStats {
	resp, err := client.Get(address + "/api/db/stats")
	require.NoError(t, err, "could not get stats")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var stats UpdateStats
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&stats), "cannot decode")
	return stats
}

type Comics struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type ComicsReply struct {
	Comics []Comics `json:"comics"`
	Total  int      `json:"total"`
}

func TestSearchNoPhrase(t *testing.T) {
	resp, err := client.Get(address + "/api/search")
	require.NoError(t, err, "failed to search")
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "need bad request")
}

func TestSearchBadLimitMinus(t *testing.T) {
	resp, err := client.Get(address + "/api/search?limit=-1")
	require.NoError(t, err, "failed to search")
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "need bad request")
}

func TestSearchBadLimitAlpha(t *testing.T) {
	resp, err := client.Get(address + "/api/search?limit=asdf")
	require.NoError(t, err, "failed to search")
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "need bad request")
}

func TestSearchLimit2(t *testing.T) {
	resp, err := client.Get(address + "/api/search?limit=2&phrase=linux")
	require.NoError(t, err, "failed to search")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "need OK status")
	var comics ComicsReply
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&comics), "decode failed")
	require.Equal(t, 2, comics.Total)
	require.Equal(t, 2, len(comics.Comics))
}

func TestSearchLimitDefault(t *testing.T) {

	resp, err := client.Get(address + "/api/search?phrase=linux")
	require.NoError(t, err, "failed to search")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "need OK status")
	var comics ComicsReply
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&comics), "decode failed")
	require.Equal(t, 10, comics.Total)
	require.Equal(t, 10, len(comics.Comics))
}

func TestSearchPhrases(t *testing.T) {
	testCases := []struct {
		phrase string
		url    string
	}{
		{
			phrase: "linux+cpu+video+machine+русские+хакеры",
			url:    "https://imgs.xkcd.com/comics/supported_features.png",
		},
		{
			phrase: "Binary Christmas Tree",
			url:    "https://imgs.xkcd.com/comics/tree.png",
		},
		{
			phrase: "apple a day -> keeps doctors away",
			url:    "https://imgs.xkcd.com/comics/an_apple_a_day.png",
		},
		{
			phrase: "mines, captcha",
			url:    "https://imgs.xkcd.com/comics/mine_captcha.png",
		},
		{
			phrase: "newton apple's idea",
			url:    "https://imgs.xkcd.com/comics/inspiration.png",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.phrase, func(t *testing.T) {
			resp, err := client.Get(address + "/api/search?phrase=" + url.QueryEscape(tc.phrase))
			require.NoError(t, err, "failed to search")
			defer resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode, "need OK status")
			var comics ComicsReply
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&comics), "decode failed")
			urls := make([]string, 0, len(comics.Comics))
			for _, c := range comics.Comics {
				urls = append(urls, c.URL)
			}
			require.Containsf(t, urls, tc.url, "could not find %q", tc.phrase)
		})
	}
}
