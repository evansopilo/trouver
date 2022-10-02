package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Review data object definition with struct tag annotation to instruct the json and bson encoder
// on how the keys of the json and bson encoded output should look like.
type Review struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	PlaceID     string    `json:"place_id,omitempty" bson:"place_id,omitempty"`
	UserID      string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	TextContent string    `json:"title,omitempty" bson:"title,omitempty"`
	Rating      float32   `json:"rating,omitempty" bson:"rating,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type Reviews []Review

type ReviewModel struct {
	client *mongo.Client
}

func NewReviewModel(client *mongo.Client) *ReviewModel { return &ReviewModel{client: client} }

// InsertOne inserts a new document to the reviews collection, takes a context, database name, collection name
// and pointer to place struct object with the data to be inserted.
func (r ReviewModel) InsertOne(ctx context.Context, database, collection string, review *Review) error {
	var result *mongo.InsertOneResult
	coll := r.client.Database(database).Collection(collection)
	result, err := coll.InsertOne(ctx, review)
	if err != nil {
		return err
	}
	if result.InsertedID.(string) != review.ID {
		return ErrNoDocument
	}
	return nil
}

// UpdateOne updates a specific review document in the reviews collection, takes a context, database name, collection name
// and pointer to review struct objet with data to be updated.
func (r ReviewModel) UpdateOne(ctx context.Context, database, collection string, review *Review) error {
	var result *mongo.UpdateResult
	coll := r.client.Database(database).Collection(collection)
	result, err := coll.UpdateOne(ctx, bson.M{"_id": review.ID}, bson.D{{Key: "$set", Value: review}})
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 || result.UpsertedCount != 1 {
		return ErrNoDocument
	}
	if result.UpsertedID.(string) != review.ID {
		return ErrNoDocument
	}
	return nil
}

// FindOne finds a specific review document in the reviews collection, takes a context, database name, collection name
// and the document id
func (r ReviewModel) FindOne(ctx context.Context, database, collection string, reviewID string) (*Review, error) {
	var result *mongo.SingleResult
	var review Review
	coll := r.client.Database(database).Collection(collection)
	result = coll.FindOne(ctx, bson.M{"_id": reviewID})
	if err := result.Decode(&review); err != nil {
		return nil, err
	}
	return &review, nil
}

// List finds all reviews documents in the reviews collections by place id, takes a context, database name
// collection name and filter.
func (r ReviewModel) List(ctx context.Context, database, collection string, placeID string, filter Filter) (*Reviews, error) {
	opts := options.Find().SetSkip(int64(filter.Skip)).SetLimit(int64(filter.Limit))
	coll := r.client.Database(database).Collection(collection)
	filterCursor, err := coll.Find(ctx, bson.M{"place_id": placeID}, opts)
	if err != nil {
		return nil, err
	}
	var reviews Reviews
	for filterCursor.Next(ctx) {
		var review Review
		if err := filterCursor.Decode(&review); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return &reviews, nil
}

// DeleteOne deletes a specific review document in the reviews collection, takes a context, database name, collection name
// and document id.
func (r ReviewModel) DeleteOne(ctx context.Context, database, collection string, reviewID string) error {
	var result *mongo.DeleteResult
	coll := r.client.Database(database).Collection(collection)
	result, err := coll.DeleteOne(ctx, bson.M{"_id": reviewID})
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return ErrNoDocument
	}
	return nil
}
