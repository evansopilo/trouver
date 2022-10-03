package data

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

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

	Review interface {
		// InsertOne inserts a new document to the reviews collection, takes a context, database name, collection name
		// and pointer to place struct object with the data to be inserted.
		InsertOne(ctx context.Context, database, collection string, review *Review) error

		// UpdateOne updates a specific review document in the reviews collection, takes a context, database name, collection name
		// and pointer to review struct objet with data to be updated.
		UpdateOne(ctx context.Context, database, collection string, review *Review) error

		// FindOne finds a specific review document in the reviews collection, takes a context, database name, collection name
		// and the document id
		FindOne(ctx context.Context, database, collection string, reviewID string) (*Review, error)

		// List finds all reviews documents in the reviews collections by place id, takes a context, database name
		// collection name and filter.
		List(ctx context.Context, database, collection string, placeID string, filter Filter) (*Reviews, error)

		// DeleteOne deletes a specific review document in the reviews collection, takes a context, database name, collection name
		// and document id.
		DeleteOne(ctx context.Context, database, collection string, reviewID string) error
	}

	Auth interface {
		// VerifyIDToken verifys a token id, takes context, firebase app and id token.
		VerifyIDToken(ctx context.Context, app *firebase.App, idToken string) (*auth.Token, error)

		// RevokeRefreshTokens revokes refresh token associated by a user account, takes context, firebase app and user uid.
		RevokeRefreshTokens(ctx context.Context, app *firebase.App, uid string) error

		// GetUser gets user by id, takes context, firebase app and user uid.
		GetUser(ctx context.Context, app *firebase.App, uid string) (*auth.UserRecord, error)

		// GetUserByEmail gets user by email, takes context, firebase app and email.
		GetUserByEmail(ctx context.Context, app *firebase.App, email string) (*auth.UserRecord, error)

		// GetUserByPhone gets user by phone, takes context, firebase app and phone.
		GetUserByPhone(ctx context.Context, app *firebase.App, phone string) (*auth.UserRecord, error)

		// CreateUser creates a new user, takes a context, firebase app and user object.
		CreateUser(ctx context.Context, app *firebase.App, user *User) (*auth.UserRecord, error)

		// UpdateUser updates an existing user, takes context, firebase app and user object.
		UpdateUser(ctx context.Context, app *firebase.App, user *User) (*auth.UserRecord, error)

		// DeleteUser delete a user by id, takes a context, firebase app and user uid.
		DeleteUser(ctx context.Context, app *firebase.App, uid string) error

		// CustomClaimsSet sets custom claims to a user, takes context firebase app, uid and claims map.
		CustomClaimsSet(ctx context.Context, app *firebase.App, uid string, claims map[string]interface{}) error
	}
}
