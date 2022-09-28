package validator

import (
	"github.com/evansopilo/trouver/internal/data"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidatePlace(place *data.Place) error {
	return validation.ValidateStruct(&place,
		validation.Field(&place.Title, validation.Required, validation.Length(0, 150)),
		validation.Field(&place.Description, validation.Required, validation.Length(0, 150)),
		validation.Field(&place.Categories, validation.Length(0, 5)),
		validation.Field(&place.ImageURL, is.URL),
		validation.Field(&place.PhoneNumber, validation.Length(0, 20)),
		validation.Field(&place.Email, is.Email),
		validation.Field(&place.Location.Address.Street1, validation.Length(0, 30)),
		validation.Field(&place.Location.Address.City, validation.Length(0, 30)),
		validation.Field(&place.Location.Address.State, validation.Length(0, 30)),
		validation.Field(&place.Location.Address.ZipCode, validation.Length(0, 30)),
		validation.Field(&place.Location.Geo.Type, validation.Length(0, 30)),
		validation.Field(&place.Location.Geo.Coordinates[0], is.Longitude),
		validation.Field(&place.Location.Geo.Coordinates[1], is.Latitude),
	)
}
