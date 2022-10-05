package data

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

// User data object definition with struct tag annotation to instruct the json
// on how the keys of the json encoded output should look like.
type User struct {
	UID           string `json:"uid,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty"`
	Password      string `json:"password,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	PhotoURL      string `json:"photo_url,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
}

type AuthModel struct {
}

// VerifyIDToken verifys a token id, takes context, firebase app and id token.
func (AuthModel) VerifyIDToken(ctx context.Context, app *firebase.App, idToken string) (*auth.Token, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// RevokeRefreshTokens revokes refresh token associated by a user account, takes context, firebase app and user uid.
func (AuthModel) RevokeRefreshTokens(ctx context.Context, app *firebase.App, uid string) error {
	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}
	if err := client.RevokeRefreshTokens(ctx, uid); err != nil {
		return err
	}
	return nil
}

// GetUser gets user by id, takes context, firebase app and user uid.
func (AuthModel) GetUser(ctx context.Context, app *firebase.App, uid string) (*auth.UserRecord, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	user, err := client.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail gets user by email, takes context, firebase app and email.
func (AuthModel) GetUserByEmail(ctx context.Context, app *firebase.App, email string) (*auth.UserRecord, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		log.Fatalf("error getting user by email %s: %v\n", email, err)
	}
	return user, nil
}

// GetUserByPhone gets user by phone, takes context, firebase app and phone.
func (AuthModel) GetUserByPhone(ctx context.Context, app *firebase.App, phone string) (*auth.UserRecord, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	user, err := client.GetUserByPhoneNumber(ctx, phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user, takes a context, firebase app and user object.
func (AuthModel) CreateUser(ctx context.Context, app *firebase.App, user *User) (*auth.UserRecord, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	params := (&auth.UserToCreate{}).
		Email(user.Email).
		EmailVerified(user.EmailVerified).
		PhoneNumber(user.PhoneNumber).
		Password(user.Password).
		DisplayName(user.DisplayName).
		PhotoURL(user.PhotoURL).
		Disabled(user.Disabled)
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateUser updates an existing user, takes context, firebase app and user object.
func (AuthModel) UpdateUser(ctx context.Context, app *firebase.App, user *User) (*auth.UserRecord, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	params := (&auth.UserToUpdate{}).
		Email(user.Email).
		EmailVerified(user.EmailVerified).
		PhoneNumber(user.PhoneNumber).
		Password(user.Password).
		DisplayName(user.DisplayName).
		PhotoURL(user.PhotoURL).
		Disabled(user.Disabled)
	u, err := client.UpdateUser(ctx, user.UID, params)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// DeleteUser delete a user by id, takes a context, firebase app and user uid.
func (AuthModel) DeleteUser(ctx context.Context, app *firebase.App, uid string) error {
	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}
	err = client.DeleteUser(ctx, uid)
	if err != nil {
		return err
	}
	return nil
}

// CustomClaimsSet sets custom claims to a user, takes context firebase app, uid and claims map.
func (AuthModel) CustomClaimsSet(ctx context.Context, app *firebase.App, uid string, claims map[string]interface{}) error {
	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}
	err = client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		return err
	}
	return err
}
