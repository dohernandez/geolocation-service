package domain

import (
	"context"
	"io"
	"sync"

	"github.com/dohernandez/geolocation-service/pkg/log"
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

// Do executes the usecase logic. To do so, the function start 2 goroutines in background to process the data
// First step, validate the data and the second persist it.
//
// Returns error if parser csv file fails otherwise:
// 		processed - hold the amount of Geolocation processed
//		accepted - hold the amount of Geolocation inserted into the db
//		discarded - hold the amount of Geolocation discarded due to invalidation or duplication
//
func (uc *importGeolocationFromCSVFileToDBUseCase) Do(ctx context.Context, f io.Reader) (processed, accepted, discarded int, err error) {
	logger := log.FromContext(ctx)

	var gs []Geolocation
	if err = gocsv.Unmarshal(f, &gs); err != nil {
		return processed, accepted, discarded, err
	}

	vch := make(chan Geolocation)
	ech := make(chan error)
	sch := make(chan Geolocation)

	var wg sync.WaitGroup
	// this along with wg.Wait() are why the error handling works and doesn't deadlock.
	finished := make(chan bool, 1)

	// Routine to validate geolocation
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, g := range gs {
			processed++
			g := g // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint

			if err := g.Validate(); err != nil {
				if logger != nil {
					logger.
						WithError(err).
						WithField("geolocation", g).
						Debugf("Error validation")
				}
				// in case of error notify that this geolocation is discarded
				ech <- err

				continue
			}
			// send the valid geolocation to the next step
			vch <- g
		}

		// Close channel to notify upstream steps, that this one has finished
		close(vch)
	}()

	// Routine to persist geolocation
	wg.Add(1)
	go func() {
		defer wg.Done()

		var wgp sync.WaitGroup

		for {
			g, ok := <-vch
			// Once the previous step close the channel, there is no more geolocation objects to process
			if !ok {
				break
			}

			wgp.Add(1)
			go func() {
				defer wgp.Done()

				g.ID = uuid.New()

				if err := uc.persister.Persist(ctx, &g); err != nil {
					if logger != nil {
						logger.
							WithError(err).
							WithField("geolocation", g).
							Debugf("Error persist")
					}
					// in case of error notify that this geolocation is discarded
					ech <- err

					return
				}

				// send the valid geolocation to the next step
				sch <- g
			}()
		}

		wgp.Wait()
	}()

	// Wait for all processes to return and then close the result chan.
	go func() {
		wg.Wait()

		close(finished)
	}()

	var fin bool
	for {
		select {
		case <-finished:
			fin = true
		case <-ech:
			discarded++
		case <-sch:
			accepted++
		}

		if fin {
			break
		}
	}

	close(ech)
	close(sch)

	return processed, accepted, discarded, nil
}
