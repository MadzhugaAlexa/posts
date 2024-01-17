package mongo

import (
	"context"
	"news/internal/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client *mongo.Client
}

func NewStorage(client *mongo.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) nextId() int32 {
	collection := s.client.Database("test").Collection("counters")

	upsert := true
	opts := &options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}

	res := collection.FindOneAndUpdate(context.Background(), bson.M{"_id": "postID"}, bson.M{"$inc": bson.M{"score": 1}}, opts)

	var i int32 = 0
	vals := &bson.D{}
	res.Decode(&vals)
	for _, d := range *vals {
		if d.Key == "score" {
			i = d.Value.(int32)
		}
	}

	return i
}

func (s *Storage) collection() *mongo.Collection {
	return s.client.Database("test").Collection("posts")
}

func (s *Storage) CreatePost(p *storage.Post) error {
	p.ID = int(s.nextId())

	_, err := s.collection().InsertOne(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FindAll() ([]storage.Post, error) {
	collection := s.collection()
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var data []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}
	return data, nil
}

func (s *Storage) Find(id int) (storage.Post, error) {
	collection := s.collection()
	filter := bson.D{{"id", id}}
	c := collection.FindOne(context.Background(), filter)

	var post storage.Post
	err := c.Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return storage.Post{}, storage.ErrNotFound
		}

		return storage.Post{}, err
	}
	return post, nil
}

func (s *Storage) DeletePost(id int) error {
	collection := s.collection()
	filter := bson.D{{"id", id}}
	dr, err := collection.DeleteOne(context.Background(), filter)

	if dr.DeletedCount == 0 {
		return storage.ErrNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdatePost(p *storage.Post) error {
	collection := s.collection()
	filter := bson.D{{"id", p.ID}}
	_, err := collection.ReplaceOne(context.Background(), filter, p)

	if err != nil {
		return err
	}

	return err
}
