package model

// Interaction ...
type Interaction struct {
	Interaction_id int    `json:"interaction_id"`
	Title          string `json:"title"`
	Article_id     string `json:"article_id"`
	Items          string `json:"items"`
}
