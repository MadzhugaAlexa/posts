package postgres

import (
	"context"
	"errors"
	"news/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreatePost(p *storage.Post) error {
	tx, _ := s.db.Begin(context.Background())

	row := tx.QueryRow(
		context.Background(),
		"INSERT INTO posts (author_id, title, content, created_at, published_at) VALUES ($1, $2, $3, $4, $5) returning id",
		p.AuthorID, p.Title, p.Content, p.CreatedAt, p.PublishedAt,
	)

	var id int
	err := row.Scan(&id)

	if err != nil {
		return err
	}

	p.ID = id

	tx.Commit(context.Background())

	return nil
}

func (s *Storage) FindAll() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(),
		"select id, author_id, title, content, created_at, published_at from posts")

	if err != nil {
		return nil, err
	}

	posts := []storage.Post{}
	for rows.Next() {
		post := storage.Post{}
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt, &post.PublishedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Storage) Find(id int) (storage.Post, error) {
	post := storage.Post{}
	row := s.db.QueryRow(context.Background(),
		"select id, author_id, title, content, created_at, published_at from posts where id = $1",
		id)

	err := row.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt, &post.PublishedAt)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (s *Storage) UpdatePost(p *storage.Post) error {
	sql := `UPDATE posts set 
		author_id = $1, title = $2, content = $3, created_at = $4, published_at = $5
        where id = $6`

	_, err := s.db.Exec(
		context.Background(),
		sql,
		p.AuthorID, p.Title, p.Content, p.CreatedAt, p.PublishedAt, p.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeletePost(id int) error {
	sql := "DELETE FROM posts WHERE id = $1"

	tag, err := s.db.Exec(context.Background(), sql, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("не найдена запись")
	}

	return nil
}
