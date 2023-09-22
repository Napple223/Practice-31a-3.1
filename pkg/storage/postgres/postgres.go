package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// абстракция на БД как объект
type Storage struct {
	db *pgxpool.Pool
}

// Функция - конструктор для создания пула подключений
func New(constr string) (*Storage, error) {
	db, err := pgxpool.New(context.Background(), constr)

	if err != nil {
		return nil, err
	}

	s := Storage{
		db: db,
	}
	return &s, nil
}

// Структура поста
type Post struct {
	Id        int
	AuthorID  int
	Title     string
	Content   string
	CreatedAt int
}

// Функция, возвращающая все существующие посты из БД
func (s *Storage) Posts() ([]Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
		id,
		author_id,
		title,
		"content",
		created_at
		FROM posts
		ORDER BY id;
	`)

	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var p Post

		err = rows.Scan(
			&p.Id,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// Функция добавляет новый пост и возвращает его id
func (s *Storage) AddPost(p Post) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO posts (title, "content") VALUES ($1, $2) RETURNING id;
	`,
		p.Title,
		p.Content,
	).Scan(&id)
	return id, err
}

// Функция для обновления поста
func (s *Storage) UpdatePost(id int, p Post) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE posts SET title=$1, "content"=$2
		WHERE id=$3;
	`,
		p.Title, p.Content, id,
	)
	return err
}

// Функция для удаления поста из БД
func (s *Storage) DeletePost(id int) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM posts
		WHERE id=$1;
	`,
		id,
	)
	return err
}
