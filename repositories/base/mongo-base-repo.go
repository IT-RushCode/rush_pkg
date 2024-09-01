package base_repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBaseRepository interface {
	GetAll(ctx context.Context, collection string, filter interface{}, limit, offset int64, results interface{}) (int64, error)
	Create(ctx context.Context, collection string, data interface{}) (interface{}, error)
	FindByID(ctx context.Context, collection string, id primitive.ObjectID, result interface{}) error
	Update(ctx context.Context, collection string, id primitive.ObjectID, update interface{}) error
	Delete(ctx context.Context, collection string, id primitive.ObjectID) error
}

type mongoBaseRepository struct {
	*mongo.Database
}

func NewMongoBaseRepository(db *mongo.Database) MongoBaseRepository {
	return &mongoBaseRepository{db}
}

func (r *mongoBaseRepository) GetAll(ctx context.Context, collection string, filter interface{}, limit, offset int64, results interface{}) (int64, error) {
	opts := options.Find().SetLimit(limit).SetSkip(offset)
	cursor, err := r.Collection(collection).Find(ctx, filter, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, results)
	if err != nil {
		return 0, err
	}
	count, err := r.Collection(collection).CountDocuments(ctx, filter)
	return count, err
}

func (r *mongoBaseRepository) Create(ctx context.Context, collection string, data interface{}) (interface{}, error) {
	res, err := r.Collection(collection).InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (r *mongoBaseRepository) FindByID(ctx context.Context, collection string, id primitive.ObjectID, result interface{}) error {
	err := r.Collection(collection).FindOne(ctx, bson.M{"_id": id}).Decode(result)
	return err
}

func (r *mongoBaseRepository) Update(ctx context.Context, collection string, id primitive.ObjectID, update interface{}) error {
	_, err := r.Collection(collection).UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *mongoBaseRepository) Delete(ctx context.Context, collection string, id primitive.ObjectID) error {
	_, err := r.Collection(collection).DeleteOne(ctx, bson.M{"_id": id})
	return err
}
