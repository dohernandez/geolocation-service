package domain

import (
	"context"
	"fmt"
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

// Do executes the usecase logic.
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

	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, g := range gs {
			processed++
			fmt.Println(processed)

			g := g // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint

			if err := g.Validate(); err != nil {
				if logger != nil {
					logger.
						WithError(err).
						WithField("geolocation", g).
						Debugf("Error validation")
				}

				ech <- err

				continue
			}

			vch <- g
		}

		close(vch)
	}()

	concurrency := 20

	wg.Add(1)
	go func() {
		defer wg.Done()

		var wgp sync.WaitGroup
		cncr := concurrency

		for {

			g, ok := <-vch
			if !ok {
				break
			}

			wgp.Add(1)
			cncr--

			go func() {
				defer func() {
					cncr++

					wgp.Done()
				}()

				g.ID = uuid.New()

				if err := uc.persister.Persist(ctx, &g); err != nil {
					if logger != nil {
						logger.
							WithError(err).
							WithField("geolocation", g).
							Debugf("Error persist")
					}

					ech <- err

					return
				}

				sch <- g
			}()

			for {
				if cncr != 0 {
					break
				}
			}
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
