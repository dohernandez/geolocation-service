package domain

import (
	"context"

	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

// Persister defines the persist api for persisting geolocation entity
type Persister interface {
	Persist(ctx context.Context, g *Geolocation) error
}

// Geolocation represent the geolocation entity
type Geolocation struct {
	ID           uuid.UUID `csv:"-" db:"id"`
	IPAddress    string    `csv:"ip_address" db:"ip_address"`
	CountryCode  string    `csv:"country_code" db:"country_code"`
	Country      string    `csv:"country" db:"country"`
	City         string    `csv:"city" db:"city"`
	Latitude     string    `csv:"latitude" db:"latitude"`
	Longitude    string    `csv:"longitude" db:"longitude"`
	MysteryValue int64     `csv:"mystery_value" db:"mystery_value"`
}

// validateCountry validates if a string is a valid country or not.
var validateCountry = validation.NewStringRule(isCountry, "must be a valid country")

// Validate validates Geolocation
//
// Validation criteria
//		IPAddress: required, is valid ip address
//		CountryCode: required, is valid country code
//		Country: required, is valid country
//		City: required
//		Latitude: required, is valid latitude
//		Longitude: required, is valid longitude
//		MysteryValue: required
//
func (g Geolocation) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.IPAddress, validation.Required, is.IP),
		validation.Field(&g.CountryCode, validation.Required, is.CountryCode2),
		validation.Field(&g.Country, validation.Required, validateCountry),
		validation.Field(&g.City, validation.Required),
		validation.Field(&g.Latitude, validation.Required, is.Latitude),
		validation.Field(&g.Longitude, validation.Required, is.Longitude),
		validation.Field(&g.MysteryValue, validation.Required),
	)
}

// isCountry checks if a string is valid country based on
// https://www.iso.org/obp/ui/#search/code/ Code Type "Officially Assigned Codes" provides thro https://github.com/asaskevich/govalidator
func isCountry(str string) bool {
	for _, entry := range govalidator.ISO3166List {
		if str == entry.EnglishShortName {
			return true
		}
	}

	return false
}
