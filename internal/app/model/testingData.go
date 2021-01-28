package model

import "testing"

// TestNews //
func TestNews(t *testing.T) *News {
	t.Helper()

	return &News{
		Title: "Пополнение линейки товаров",
		Img:   "testImg",
	}
}
