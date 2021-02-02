package apiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println(os.Getenv("DB_PORT"))
	return func(w http.ResponseWriter, r *http.Request) {
		p := []*model.Post{}
		h := []*model.Header{}
		a := []*model.Article{}
		i := []*model.InnerDescription{}

		p, h, a, i, err := s.store.Post().GetAllPosts()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		log.Println(p)
		log.Println(h)
		log.Println(a)
		log.Println(i)
		//s.respond(w, r, http.StatusCreated, p, h, a, i)
	}
}

func (s *server) handleGetPostByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &model.Post{}
		h := &model.Header{}
		a := &model.Article{}
		i := &model.InnerDescription{}

		r.ParseForm()

		id, err := strconv.Atoi(r.Form["Post_id"][0])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err) // Спросить по поводу паники сервера...
			return
		}

		p, h, a, i, err = s.store.Post().GetPostByID(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		log.Println(p)
		log.Println(h)
		log.Println(a)
		log.Println(i)
		//s.respond(w, r, http.StatusCreated, p, h, a, i)
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
		r.ParseForm()

		id, err := strconv.Atoi(r.Form["News_id"][0])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		n := &model.News{}

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
