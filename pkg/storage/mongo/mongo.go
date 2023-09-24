package mongo

import (
	"context"

	"Practice-31a-3.1/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName   = "GoNews"
	collName = "posts"
)

type Storage struct {
	db *mongo.Client
}

func New(constr string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}

	//defer db.Disconnect(context.Background())
	err = db.Ping(context.Background(), nil)

	if err != nil {
		return nil, err
	}

	s := Storage{
		db: db,
	}

	return &s, err
}

// Функция для закрытия соединения с БД
func CloseConn(s *Storage) error {
	err := s.db.Disconnect(context.Background())
	return err
}

// Функция для получения всех постов
func (s *Storage) Posts() ([]storage.Post, error) {
	coll := s.db.Database(dbName).Collection(collName)
	filter := bson.D{}
	cur, err := coll.Find(context.Background(), filter)

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	var posts []storage.Post

	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)

		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, cur.Err()
}

// Функция для добавления поста
func (s *Storage) AddPost(p storage.Post) error {
	coll := s.db.Database(dbName).Collection(collName)
	_, err := coll.InsertOne(context.Background(), p)
	return err
}

// Функция для обновления поста
func (s *Storage) UpdatePost(p storage.Post) error {
	coll := s.db.Database(dbName).Collection(collName)
	filter := bson.D{{Key: "ID", Value: p.ID}}
	update := bson.D{{Key: "Title", Value: p.Title}, {Key: "Content", Value: p.Content}}
	_, err := coll.UpdateOne(context.Background(), filter, update)
	return err
}

// Функция для удаления поста
func (s *Storage) DeletePost(p storage.Post) error {
	coll := s.db.Database(dbName).Collection(collName)
	filter := bson.D{{Key: "ID", Value: p.ID}}
	_, err := coll.DeleteOne(context.Background(), filter)
	return err
}
