package model

// News ...
type News struct {
	News_id    int    `json:"news_id"`
	Title      string `json:"title"`
	Img        string `json:"img"`
	Date       string `json:"date"`
	Views      int    `json:"views"`
	Created_at string `json:"created_at"`
}
