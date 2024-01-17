package storage

import "errors"

type Post struct {
	ID          int
	Title       string
	Content     string
	AuthorID    int
	AuthorName  string
	CreatedAt   int64
	PublishedAt int64
}

var ErrNotFound = errors.New("not found")
