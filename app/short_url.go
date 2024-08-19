package app

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShortURL struct {
	ID       string     `bson:"_id"`
	URL      string     `bson:"url"`
	ExpireAt *time.Time `bson:"expireAt,omitempty"`
}

type UrlDAO struct {
	collection *mongo.Collection
}

func NewUrlDAO(ctx context.Context, client *mongo.Client) (*UrlDAO, error) {
	dao := &UrlDAO{
		collection: client.Database("mongo").Collection("shortUrls"),
	}

	if err := dao.create(ctx); err != nil {
		return nil, err
	}

	return dao, nil
}

func (dao *UrlDAO) create(ctx context.Context) error {
	_, err := dao.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expireAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	return err
}

func (dao *UrlDAO) Insert(ctx context.Context, shortURL *ShortURL) error {
	_, err := dao.collection.InsertOne(ctx, shortURL)
	return err
}

func (dao *UrlDAO) FindByID(ctx context.Context, id string) (*ShortURL, error) {
	var shortURL ShortURL
	filter := bson.D{{Key: "_id", Value: id}}
	err := dao.collection.FindOne(ctx, filter).Decode(&shortURL)
	switch err {
	case nil:
		return &shortURL, nil
	case mongo.ErrNoDocuments:
		return nil, errors.Errorf("not found")
	default:
		return nil, err
	}
}

func (dao *UrlDAO) DeleteByID(ctx context.Context, id string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := dao.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.Errorf("not found")
	}
	return nil
}
