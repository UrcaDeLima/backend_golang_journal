package store

import "github.com/UrcaDeLima/backend_golang_journal/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}

// NewsRepository ...
type NewsRepository interface {
	CreateNews(*model.News) error
}
