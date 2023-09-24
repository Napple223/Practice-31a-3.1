package postgres

import (
	"context"

	"Practice-31a-3.1/pkg/storage"
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
// type Post struct {
// 	Id        int
// 	AuthorID  int
// 	Title     string
// 	Content   string
// 	CreatedAt int
// }

// Функция, возвращающая все существующие посты из БД
func (s *Storage) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
		posts.id,
		posts.title,
		posts."content",
		posts.author_id,
		authors.name
		posts.created_at
		FROM posts, authors
		WHERE posts.author_id=authors.id
		ORDER BY posts.id;
	`)

	if err != nil {
		return nil, err
	}

	var posts []storage.Post

	for rows.Next() {
		var p storage.Post

		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// Функция добавляет новый пост
func (s *Storage) AddPost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts (title, "content") VALUES ($1, $2);
	`,
		p.Title,
		p.Content,
	)
	return err
}

// Функция для обновления поста
func (s *Storage) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE posts SET title=$1, "content"=$2
		WHERE id=$3;
	`,
		p.Title, p.Content, p.ID,
	)
	return err
}

// Функция для удаления поста
func (s *Storage) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM posts
		WHERE id=$1;
	`,
		p.ID,
	)
	return err
}
