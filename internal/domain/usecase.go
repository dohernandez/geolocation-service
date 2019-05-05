package domain

import (
	"context"
	"io"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

type (
	// ImportGeolocationFromCSVFileUseCase defines use case api for importing geolocation data
	// from csv file
	ImportGeolocationFromCSVFileUseCase interface {
		Do(ctx context.Context, f io.Reader) (processed, accepted, discarded int, err error)
	}
)

type importGeolocationFromCSVFileToDBUseCase struct {
	persister Persister
}

// NewImportGeolocationFromCSVFileToDBUseCase create an instance of ImportGeolocationFromCSVFileUseCase
// that will import the geolocation data from csv file into the database
func NewImportGeolocationFromCSVFileToDBUseCase(persister Persister) ImportGeolocationFromCSVFileUseCase {
	return &importGeolocationFromCSVFileToDBUseCase{
		persister: persister,
	}
}

// Do executes the usecase logic.
//
// Returns error if parser csv file fails otherwise:
// 		processed - hold the amount of Geolocation processed
//		accepted - hold the amount of Geolocation inserted into the db
//		discarded - hold the amount of Geolocation discarded due to invalidation or duplication
//
func (uc *importGeolocationFromCSVFileToDBUseCase) Do(ctx context.Context, f io.Reader) (processed, accepted, discarded int, err error) {
	var gs []Geolocation
	if err = gocsv.Unmarshal(f, &gs); err != nil {
		return processed, accepted, discarded, err
	}

	for _, g := range gs {
		processed++

		if err := g.Validate(); err != nil {
			discarded++

			continue
		}

		g.ID = uuid.New()

		// Using a reference for the variable on range scope `g` (scopelint)
		pg := g
		if err := uc.persister.Persist(ctx, &pg); err != nil {
			discarded++

			continue
		}

		accepted++
	}

	return processed, accepted, discarded, nil
}
