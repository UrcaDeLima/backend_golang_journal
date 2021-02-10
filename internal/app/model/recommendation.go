package model

// Recommendation ...
type Recommendation struct {
	Recommendation_id int    `json:"recommendation_id"`
	Title             string `json:"title"`
	Article_id        string `json:"article_id"`
	Text              string `json:"text"`
}
