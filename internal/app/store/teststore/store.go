package teststore

// import (
// 	"github.com/UrcaDeLima/backend_golang_journal/internal/app/model"
// 	"github.com/UrcaDeLima/backend_golang_journal/internal/app/store"
// )

// // Store ...
// type Store struct {
// 	userRepository *UserRepository
// }

// // New ...
// func New() *Store {
// 	return &Store{}
// }

// // User ...
// func (s *Store) User() store.UserRepository {
// 	if s.userRepository != nil {
// 		return s.userRepository
// 	}

// 	s.userRepository = &UserRepository{
// 		store: s,
// 		users: make(map[string]*model.User),
// 	}

// 	return s.userRepository
// }
