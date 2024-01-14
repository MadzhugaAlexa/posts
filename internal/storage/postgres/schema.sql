DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES authors(id) NOT NULL,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    published_at BIGINT NOT NULL
);

INSERT INTO authors (name) VALUES ('Александра');
INSERT INTO posts (author_id, title, content, created_at, published_at) VALUES (1, 'Статья', 'Содержание статьи', 0, 0);
