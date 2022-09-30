package data

import "context"

type Models struct {
	Place interface {
		// InsertOne inserts a new document to the places collection, takes a context, database name, collection name
		// and pointer to place struct object with the data to be inserted.
		InsertOne(ctx context.Context, database, collection string, place *Place) error

		// UpdateOne updated a specific place document in the places collection, takes a context, database name, collection name
		// and pointer to place struct objet with data to be updated.
		UpdateOne(ctx context.Context, database, collection string, place *Place) error

		// FindOne finds a specific places document in the places collection, takes a context, database name, collection name
		// and the document id
		FindOne(ctx context.Context, database, collection string, placeID string) (*Place, error)

		// List finds all places documents in the places collections, takes a context, database name, collection name
		// and filter.
		List(ctx context.Context, database, collection string, filter Filter) (*Places, error)

		// DeleteOne deletes a specific place document in the places collection, takes a context, database name, collection name
		// and document id.
		DeleteOne(ctx context.Context, database, collection string, placeID string) error

		// SearchPlace searches place documents in places collection by search term, takes a context, database name, collection name
		// search term and filter.
		SearchPlace(ctx context.Context, database, collection string, term string, filter Filter) (*Places, error)
	}
}
