package apiserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/UrcaDeLima/backend_golang_journal/internal/app/model"
	"github.com/UrcaDeLima/backend_golang_journal/internal/app/store"
	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/setNews", s.handleCreateNews()).Methods("POST")
	s.router.HandleFunc("/getNewsById", s.handleGetNewsByID()).Methods("GET")
	s.router.HandleFunc("/getAllNews", s.handleGetAllNews()).Methods("GET")
	s.router.HandleFunc("/getPostById", s.handleGetPostByID()).Methods("GET")
	s.router.HandleFunc("/getAllPosts", s.handleGetAllPosts()).Methods("GET")
}

func (s *server) handleGetAllPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postModel := []*model.PostModel{}

		postModel, err := s.store.Post().GetAllPosts()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, postModel)
	}
}

func (s *server) handleGetPostByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postModel := &model.PostModel{}

		r.ParseForm()

		id, err := strconv.Atoi(r.Form["id"][0]) // Тут сервер падает, при неверном ключе...
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		postModel, err = s.store.Post().GetPostByID(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, postModel)
	}
}

func (s *server) handleGetAllNews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := []*model.News{}

		n, err := s.store.News().GetAllNews()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, n)
	}
}

func (s *server) handleGetNewsByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := &model.News{}

		r.ParseForm()

		id, err := strconv.Atoi(r.Form["id"][0]) // Тут сервер падает, при неверном ключе...
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		n, err = s.store.News().GetNewsByID(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, n)
	}
}

func (s *server) handleCreateNews() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
		Img   string `json:"img"`
		Views int    `json:"views"`
		Date  string `json:"date"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		n := &model.News{
			Title: req.Title,
			Img:   req.Img,
			Views: req.Views,
			Date:  req.Date,
		}

		if err := s.store.News().CreateNews(n); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, n)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
