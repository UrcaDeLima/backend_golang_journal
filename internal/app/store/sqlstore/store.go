package sqlstore

import (
	"database/sql"

	"github.com/UrcaDeLima/backend_golang_journal/internal/app/store"
)

// Store ...
type Store struct {
	db                         *sql.DB
	newsRepository             *NewsRepository
	headerRepository           *HeaderRepository
	articleRepository          *ArticleRepository
	postRepository             *PostRepository
	innerDescriptionRepository *InnerDescriptionRepository
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// News ...
func (s *Store) News() store.NewsRepository {
	if s.newsRepository != nil {
		return s.newsRepository
	}

	s.newsRepository = &NewsRepository{
		store: s,
	}

	return s.newsRepository
}

// Header ...
func (s *Store) Header() store.HeaderRepository {
	if s.headerRepository != nil {
		return s.headerRepository
	}

	s.headerRepository = &HeaderRepository{
		store: s,
	}

	return s.headerRepository
}

// Article ...
func (s *Store) Article() store.ArticleRepository {
	if s.articleRepository != nil {
		return s.articleRepository
	}

	s.articleRepository = &ArticleRepository{
		store: s,
	}

	return s.articleRepository
}

// Post ...
func (s *Store) Post() store.PostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}

	s.postRepository = &PostRepository{
		store: s,
	}

	return s.postRepository
}

// InnerDescription ...
func (s *Store) InnerDescription() store.InnerDescriptionRepository {
	if s.innerDescriptionRepository != nil {
		return s.innerDescriptionRepository
	}

	s.innerDescriptionRepository = &InnerDescriptionRepository{
		store: s,
	}

	return s.innerDescriptionRepository
}
