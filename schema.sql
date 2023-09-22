/* Схема БД для PostgresSQL */

DROP TABLE IF EXISTS authors, posts;

CREATE TABLE IF NOT EXISTS authors (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
	id SERIAL PRIMARY KEY,
	author_id INTEGER REFERENCES authors(id) DEFAULT 0,
	title TEXT NOT NULL,
	"content" TEXT NOT NULL,
	created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch from now())
);

/*default*/

INSERT INTO authors (id, name) VALUES (0, 'default')
