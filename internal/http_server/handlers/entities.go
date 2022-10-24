package handlers

type ErrorResp struct {
	Error string `json:"error"`
}

type NewEntity struct {
	Id      uint64 `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"description"`
	Link    string `json:"link"`
	PubDate string `json:"pub_date"`
}

type GetNewsResp struct {
	News []NewEntity `json:"news"`
}

type GetNewResp struct {
	Id      uint64 `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"description"`
	Link    string `json:"link"`
	PubDate string `json:"pub_date"`
}
