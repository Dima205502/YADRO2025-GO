package model

type Comics struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type ComicsResponse struct {
	Comics []Comics `json:"comics"`
	Total  int      `json:"total"`
}

type TemplateData struct {
	Phrase       string
	SearchID     string
	Comics       []Comics
	Total        int
	CurrentIndex int
	DisplayTotal int
}

type AuthInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Status struct {
	Status string `json:"status"`
}

type StatsResponse struct {
	WordsTotal    int `json:"words_total"`
	WordsUnique   int `json:"words_unique"`
	ComicsFetched int `json:"comics_fetched"`
	ComicsTotal   int `json:"comics_total"`
}
