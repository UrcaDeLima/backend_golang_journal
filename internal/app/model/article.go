package model

// Article ...
type Article struct {
	Article_id    int      `json:"article_id"`
	Title         string   `json:"title"`
	BackgroundImg string   `json:"backgroundImg"`
	Paragraphs    []string `json:"paragraphs"`
	Text          string   `json:"text"`
	Post_id       int      `json:"post_id"`
	Created_at    string   `json:"created_at"`
}
