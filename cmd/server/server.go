package main

import (
	"context"
	"log"
	"net/http"
	"news/internal/api"
	"news/internal/storage/postgres"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	DB_URL := "postgres://alexa:alexa@localhost:5432/posts"

	db, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v\n", err)
	}
	defer db.Close()

	st := postgres.NewStorage(db)

	api := api.NewStorage(st)

	mux := http.NewServeMux()

	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idParam := strings.TrimPrefix(r.URL.Path, "/posts/")

		switch r.Method {
		case http.MethodGet:
			if idParam == "" {
				api.FindAll(w, r)
			} else {
				api.FindOne(idParam, w, r)
			}
		case http.MethodPost:
			api.CreatePost(w, r)
		case http.MethodPut, http.MethodPatch:
			api.Update(idParam, w, r)
		case http.MethodDelete:
			api.Delete(idParam, w, r)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
