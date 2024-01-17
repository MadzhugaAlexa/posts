package main

import (
	"context"
	"log"
	"net/http"
	"news/internal/api"
	"strings"

	m "news/internal/storage/mongo"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	DB_URL := "postgres://alexa:alexa@localhost:5432/posts"

	pg, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v\n", err)
	}
	defer pg.Close()

	mongoOpts := options.Client().ApplyURI("mongodb://localhost:27017/")
	mn, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer mn.Disconnect(context.Background())

	err = mn.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// pgSt := postgres.NewStorage(pg)
	mgSt := m.NewStorage(mn)
	// memSt := memdb.NewStorage()

	// api := api.NewStorage(memSt)
	api := api.NewStorage(mgSt)
	// api := api.NewStorage(st)

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
