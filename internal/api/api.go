package api

import (
	"encoding/json"
	"net/http"
	"news/internal/storage"
	"strconv"
)

type Storage interface {
	CreatePost(p *storage.Post) error
	FindAll() ([]storage.Post, error)
	Find(id int) (storage.Post, error)
	UpdatePost(p *storage.Post) error
	DeletePost(id int) error
}

type Api struct {
	storage Storage
}

func NewStorage(st Storage) *Api {
	return &Api{
		storage: st,
	}
}

func (a *Api) FindAll(w http.ResponseWriter, r *http.Request) {
	posts, err := a.storage.FindAll()
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(posts)
}

func (a *Api) FindOne(idParam string, w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(idParam)

	post, err := a.storage.Find(id)
	if err != nil {
		if err == storage.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(post)
}

func (a *Api) CreatePost(w http.ResponseWriter, r *http.Request) {
	post := &storage.Post{}
	err := json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		panic(err)
	}
	err = a.storage.CreatePost(post)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(post)
}

func (a *Api) Update(idParam string, w http.ResponseWriter, r *http.Request) {
	post := storage.Post{}
	json.NewDecoder(r.Body).Decode(&post)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		panic(err)
	}

	post.ID = id

	err = a.storage.UpdatePost(&post)
	if err != nil {
		if err == storage.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(post)
}

func (a *Api) Delete(idParam string, w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(idParam)

	err := a.storage.DeletePost(id)
	if err != nil {
		if err == storage.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	return
}
