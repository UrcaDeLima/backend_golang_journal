package model

// Post ...
type Post struct {
	Post_id             int    `json:"post_id"`
	InnerAdvertising_id string `json:"innerAdvertising_id"`
	Created_at          string `json:"created_at"`
}
