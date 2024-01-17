package memdb

import (
	"news/internal/storage"
	"sync"
)

type Storage struct {
	posts map[int]*storage.Post
	mu    sync.RWMutex
	id    int
}

func NewStorage() *Storage {
	return &Storage{
		posts: make(map[int]*storage.Post),
	}
}

func (s *Storage) CreatePost(p *storage.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id = s.id + 1

	p.ID = s.id
	s.posts[p.ID] = p
	return nil
}

func (s *Storage) FindAll() ([]storage.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := []storage.Post{}

	for _, p := range s.posts {
		posts = append(posts, *p)
	}

	return posts, nil
}

func (s *Storage) Find(id int) (storage.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post := s.posts[id]
	if post == nil {
		return storage.Post{}, storage.ErrNotFound
	}

	return *post, nil
}

func (s *Storage) UpdatePost(p *storage.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.posts[p.ID] == nil {
		return storage.ErrNotFound
	}

	s.posts[p.ID] = p

	return nil
}

func (s *Storage) DeletePost(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.posts[id] == nil {
		return storage.ErrNotFound
	}

	delete(s.posts, id)
	return nil
}
