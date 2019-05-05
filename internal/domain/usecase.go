package domain

import (
	"context"
	"io"

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
// Returns error if parser csv file fails.
//
// When log is provide thro ctx, execution details are logged.
//
// Example of usage:
//
//		// ... file *os.File variable is initialized
//
//		uc := domain.NewImportGeolocationFromCSVFileToDBUseCase()
//
//		// setting log into context to allow to log the execution details
// 		l := logrus.New()
// 		ctx := log.ToContext(context.TODO(), l)
//
//		uc.Do(ctx, file)
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

	//// getting csv header
	//r := csv.NewReader(f)
	//line, err := r.Read()
	//if err != nil {
	//	return err
	//}
	//
	//var processed, accepted, discarded int
	//
	//header := strings.Join(line, ",")
	//
	//// parsing csv lines.
	//// The reason why the file is parser line by line and not all at once is because for those lines where
	//// the value is not valid type an error is return, therefore the function fails without parse any line.
	//for {
	//	line, err := r.Read()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		return err
	//	}
	//
	//	processed++
	//
	//	var gs []Geolocation
	//
	//	// concatenate header with line to allow unmarshall into Geolocation struct
	//	s := header + "\n" + strings.Join(line, ",")
	//
	//	// creating inline file
	//	lr := strings.NewReader(s)
	//	if err := gocsv.Unmarshal(lr, &gs); err != nil {
	//		discarded++
	//
	//		continue
	//	}
	//
	//	g := gs[0]
	//	if err := g.Validate(); err != nil {
	//		discarded++
	//
	//		continue
	//	}
	//
	//	accepted++
	//}

	l := log.FromContext(ctx)
	if l != nil {
		l.
			WithField("processed", processed).
			WithField("accepted", accepted).
			WithField("discarded", discarded).
			Infof("Import Geolocation From CSV to DB")
	}

	return processed, accepted, discarded, nil
}
