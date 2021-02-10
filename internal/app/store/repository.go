package store

import (
	"mime/multipart"

	"github.com/UrcaDeLima/backend_golang_journal/internal/app/model"
)

// NewsRepository ...
type NewsRepository interface {
	CreateNews(*model.News) error
	GetNewsByID(id int) (*model.News, error)
	GetAllNews() ([]*model.News, error)
}

// HeaderRepository ...
type HeaderRepository interface {
}

// ArticleRepository ...
type ArticleRepository interface {
}

// InnerDescriptionRepository ...
type InnerDescriptionRepository interface {
}

// InteractionRepository ...
type InteractionRepository interface {
}

// RecommendationRepository ...
type RecommendationRepository interface {
}

// ImageRepository ...
type ImageRepository interface {
	SetPicture(m *multipart.Reader) error
	UpdatePicture(id int, m *multipart.Reader) error
}

// PostRepository ...
type PostRepository interface {
	GetPostByID(id int) (*model.PostModel, error)
	GetAllPosts() ([]*model.PostModel, error)
	CreatePost(m *multipart.Reader)
}
