package repository

import (
	"context"
	"errors"
	"time"

	"github.com/hussainr95/link-in-bio-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkRepository interface {
	Create(ctx context.Context, link *entity.Link) (*entity.Link, error)
	GetByID(ctx context.Context, id string) (*entity.Link, error)
	Update(ctx context.Context, link *entity.Link) (*entity.Link, error)
	Delete(ctx context.Context, id string) error
	IncrementClicks(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

type mongoLinkRepository struct {
	collection *mongo.Collection
}

func NewMongoLinkRepository(db *mongo.Database) LinkRepository {
	return &mongoLinkRepository{
		collection: db.Collection("links"),
	}
}

func (r *mongoLinkRepository) Create(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	res, err := r.collection.InsertOne(ctx, link)
	if err != nil {
		return nil, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		link.ID = oid.Hex()
	}
	return link, nil
}

func (r *mongoLinkRepository) GetByID(ctx context.Context, id string) (*entity.Link, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var link entity.Link
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *mongoLinkRepository) Update(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	if link.ID == "" {
		return nil, errors.New("missing link ID")
	}

	oid, err := primitive.ObjectIDFromHex(link.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"title":     link.Title,
			"url":       link.URL,
			"expiresAt": link.ExpiresAt,
			"clicks":    link.Clicks,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	// Optionally re-fetch the updated document
	return r.GetByID(ctx, link.ID)
}

func (r *mongoLinkRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *mongoLinkRepository) IncrementClicks(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$inc": bson.M{"clicks": 1}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoLinkRepository) DeleteExpired(ctx context.Context) error {
	// Delete all links with expiresAt before now
	_, err := r.collection.DeleteMany(ctx, bson.M{"expiresAt": bson.M{"$lt": time.Now()}})
	return err
}
