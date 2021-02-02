package model

// Header ...
type Header struct {
	Header_id        int    `json:"header_id"`
	Title            string `json:"title"`
	Image_id         int    `json:"image_id"`
	Date             string `json:"date"`
	Views            int    `json:"views"`
	ShortDescription string `json:"shortDescription"`
	Post_id          int    `json:"post_id"`
	Created_at       string `json:"created_at"`
}
