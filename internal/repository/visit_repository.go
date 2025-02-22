package repository

import (
	"context"

	"github.com/hussainr95/link-in-bio-service/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type VisitRepository interface {
	Create(ctx context.Context, visit *entity.Visit) (*entity.Visit, error)
}

type mongoVisitRepository struct {
	collection *mongo.Collection
}

func NewMongoVisitRepository(db *mongo.Database) VisitRepository {
	return &mongoVisitRepository{
		collection: db.Collection("visits"),
	}
}

func (r *mongoVisitRepository) Create(ctx context.Context, visit *entity.Visit) (*entity.Visit, error) {
	_, err := r.collection.InsertOne(ctx, visit)
	if err != nil {
		return nil, err
	}
	return visit, nil
}
