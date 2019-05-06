package domain_test

import (
	"fmt"
	"testing"

	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestValidateGeolocation(t *testing.T) {
	g := domain.Geolocation{
		IPAddress:    "160.103.7.140",
		CountryCode:  "CZ",
		Country:      "Nicaragua",
		City:         "New Neva",
		Latitude:     "-68.31023296602508",
		Longitude:    "-37.62435199624531",
		MysteryValue: 7301823115,
	}

	testCases := []struct {
		scenario string
		assert   func(t *testing.T, g domain.Geolocation)
	}{
		{
			scenario: "Invalid IP address, it is empty",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting IPAddress to empty value
				g.IPAddress = ""

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid IP address, it is not IP string",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting IPAddress to not valid ip address string
				g.IPAddress = "ip"

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid country code, it is empty",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting CountryCode to empty value
				g.CountryCode = ""

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid country code, it is not country code string",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting Country to not valid country code string
				g.CountryCode = "code"

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid country, it is empty",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting Country to empty value
				g.Country = ""

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid country, it is not country string",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting Country to not valid country string
				g.Country = "country"

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid city, it is empty",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting City to not valid country string
				g.City = ""

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid latitude, it is empty",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting Latitude to empty value
				g.Latitude = ""

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid longitude, it is not latitude string",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting Latitude to not valid latitude string
				g.Latitude = "latitude"

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid longitude, it is empty",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting Longitude to empty value
				g.Longitude = ""

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid longitude, it is not longitude string",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting Longitude to not valid longitude string
				g.Longitude = "longitude"

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Invalid mystery value, it is empty",
			assert: func(t *testing.T, g domain.Geolocation) {
				// Setting MysteryValue to empty value
				g.MysteryValue = 0

				err := g.Validate()
				fmt.Println(err)
				assert.Error(t, err)
			},
		},
		{
			scenario: "Valid",
			assert: func(t *testing.T, g domain.Geolocation) {
				err := g.Validate()
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint
		t.Run(tc.scenario, func(t *testing.T) {
			tc.assert(t, g)
		})
	}
}
