package apiserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	s.router.HandleFunc("/setPicture", s.handleSetPicture()).Methods("POST")
	s.router.HandleFunc("/getNewsById/{id}", s.handleGetNewsByID()).Methods("GET")
	s.router.HandleFunc("/getAllNews", s.handleGetAllNews()).Methods("GET")
	s.router.HandleFunc("/getPostById", s.handleGetPostByID()).Methods("GET")
	s.router.HandleFunc("/getAllPosts", s.handleGetAllPosts()).Methods("GET")
}

func (s *server) handleSetPicture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile("img") // img is the key of the form-data
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		fmt.Println("File is good")
		fmt.Println(handler.Filename)
		fmt.Println()
		fmt.Println(handler.Header)

		f, err := os.OpenFile("image/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		// srcPath := handler.Filename
		// dstPath := "image/" + handler.Filename
		// err = os.Rename(srcPath, dstPath)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// log.Println(r)
		// req := &request{}
		// if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		// 	log.Println(req)
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// 	return
		// }
		// log.Println(req)

		// err := req.ParseForm()
		// if err != nil {
		// 	// in case of any error
		// 	return
		// }
		// value := req.PostForm
		//vars := mux.Vars(req)
		//vars := req.URL.Query()
		//itemID := vars["id"]
		//assetID := newAssetID()
		//var test string

		//log.Println(value)
		//vars := req.URL.Query()

		// id, err := strconv.Atoi(vars["id"][0]) // Тут сервер падает, при неверном ключе...
		// if err != nil {
		// 	return
		// }

		// log.Println(id)
		//log.Println(vars)
		// if err := json.NewDecoder(req.Body).Decode(test); err != nil {
		// 	return
		// }
		// log.Printf("%+q\n", test)
		// log.Println(3333)
		// file, _, err := req.FormFile("image")
		// if err != nil {
		// 	log.Print(err)
		// 	log.Print("1")
		// 	resp.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		// log.Print("2")
		// // verify image
		// img, _, err := image.Decode(file)
		// if err != nil {
		// 	log.Printf("could not decode body into an image")
		// 	resp.Header().Add("Access-Control-Allow-Origin", "*")
		// 	resp.WriteHeader(http.StatusBadRequest)
		// 	resp.Write([]byte("could not decode body image"))
		// 	return
		// }

		// log.Print("3")
		// s.store.Post().SetPicture(img)

		// log.Print("4")
		//log.Println(imgRes)
	}
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
		vars := r.URL.Query()
		//fmt.Println(vars["id"][0])
		n := &model.News{}
		//log.Println(vars["id"])
		//r.ParseForm()

		id, err := strconv.Atoi(vars["id"][0]) // Тут сервер падает, при неверном ключе...
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
