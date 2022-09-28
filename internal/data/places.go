package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Place data object definition with struct tag annotation to instruct the json and bson encoder
// on how the keys of the json and bson encoded output should look like.
type Place struct {
	ID          string   `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Title       string   `json:"title,omitempty" bson:"title,omitempty"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Categories  []string `json:"categories,omitempty" bson:"categories,omitempty"`
	ImageURL    string   `json:"image_url,omitempty" bson:"image_url,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Email       string   `json:"email,omitempty" bson:"email,omitempty"`
	Location    struct {
		Address struct {
			Street1 string `json:"street_1,omitempty" bson:"street_1,omitempty"`
			City    string `json:"city,omitempty" bson:"city,omitempty"`
			State   string `json:"state,omitempty" bson:"state,omitempty"`
			ZipCode string `json:"zip_code,omitempty" bson:"zip_code,omitempty"`
		} `json:"address,omitempty" bson:"address,omitempty"`
		Geo struct {
			Type        string    `json:"type,omitempty" bson:"type,omitempty"`
			Coordinates []float64 `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
		} `json:"geo,omitempty" bson:"geo,omitempty"`
	} `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type Places []Place

type PlaceModel struct {
	client *mongo.Client
}

func NewPlaceModel(client *mongo.Client) *PlaceModel { return &PlaceModel{client: client} }

// InsertOne inserts a new document to the places collection, takes a context, database name, collection name
// and pointer to place struct object with the data to be inserted.
func (p PlaceModel) InsertOne(ctx context.Context, database, collection string, place *Place) error {
	var result *mongo.InsertOneResult
	coll := p.client.Database(database).Collection(collection)
	result, err := coll.InsertOne(ctx, place)
	if err != nil {
		return err
	}
	if result.InsertedID.(string) != place.ID {
		return err
	}
	return nil
}

// UpdateOne updated a specific place document in the places collection, takes a context, database name, collection name
// and pointer to place struct objet with data to be updated.
func (p PlaceModel) UpdateOne(ctx context.Context, database, collection string, place *Place) error {
	var result *mongo.UpdateResult
	coll := p.client.Database(database).Collection(collection)
	result, err := coll.UpdateOne(ctx, bson.M{"_id": place.ID}, bson.D{{Key: "$set", Value: place}})
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 || result.UpsertedCount != 1 {
		return err
	}
	if result.UpsertedID.(string) != place.ID {
		return err
	}
	return nil
}

// FindOne finds a specific places document in the places collection, takes a context, database name, collection name
// and the document id
func (p PlaceModel) FindOne(ctx context.Context, database, collection string, placeID string) (*Place, error) {
	var result *mongo.SingleResult
	var place Place
	coll := p.client.Database(database).Collection(collection)
	result = coll.FindOne(ctx, bson.M{"_id": placeID})
	if err := result.Decode(&place); err != nil {
		return nil, err
	}
	return &place, nil
}

// List finds all places documents in the places collections, takes a context, database name, collection name
// and filter.
func (p PlaceModel) List(ctx context.Context, database, collection string, filter Filter) (*Places, error) {
	opts := options.Find().SetSkip(int64(filter.Skip)).SetLimit(int64(filter.Limit))
	coll := p.client.Database(database).Collection(database)
	filterCursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var places Places
	for filterCursor.Next(ctx) {
		var place Place
		if err := filterCursor.Decode(&place); err != nil {
			return nil, err
		}
		places = append(places, place)
	}
	return &places, nil
}

// DeleteOne deletes a specific place document in the places collection, takes a context, database name, collection name
// and document id.
func (p PlaceModel) DeleteOne(ctx context.Context, database, collection string, placeID string) error {
	var result *mongo.DeleteResult
	coll := p.client.Database(database).Collection(collection)
	result, err := coll.DeleteOne(ctx, bson.M{"_id": placeID})
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return err
	}
	return nil
}

// SearchPlace searches place documents in places collection by search term, takes a context, database name, collection name
// search term and filter.
func (p PlaceModel) SearchPlace(ctx context.Context, database, collection string, term string, filter Filter) (*Places, error) {
	sort := bson.D{{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}}}
	opts := options.Find().SetSkip(int64(filter.Skip)).SetLimit(int64(filter.Limit)).SetSort(sort)
	filt := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: term}}}}
	coll := p.client.Database(database).Collection(collection)
	filterCursor, err := coll.Find(ctx, filt, opts)
	if err != nil {
		return nil, err
	}
	var places Places
	for filterCursor.Next(ctx) {
		var place Place
		if err := filterCursor.Decode(&place); err != nil {
			return nil, err
		}
		places = append(places, place)
	}
	return &places, nil
}
