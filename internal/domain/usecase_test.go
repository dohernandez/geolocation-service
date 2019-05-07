// +build !race

package domain_test

import (
	"context"
	"strings"
	"testing"

	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/pkg/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestImportGeolocationFromCSVFileUseCase(t *testing.T) {
	testCases := []struct {
		scenario      string
		persisterFunc func(g *domain.Geolocation) error
		assert        func(t *testing.T, uc domain.ImportGeolocationFromCSVFileUseCase)
	}{
		{
			scenario: "Import data successfully",
			persisterFunc: func(g *domain.Geolocation) error {
				// checking that the geolocation object is valid
				assert.NoError(t, g.Validate())
				// checking that the geolocation object has ID
				assert.NotEqual(t, uuid.Nil, g.ID)

				return nil
			},
			assert: func(t *testing.T, uc domain.ImportGeolocationFromCSVFileUseCase) {
				in := `ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
17.78.52.164,PG,New Caledonia,Beckerberg,,,0
152.211.161.240,EU,Armenia,New Kennithbury,latitude,longitude,0`

				l := logrus.New()
				ctx := log.ToContext(context.TODO(), l)

				processed, accepted, discarded, err := uc.Do(ctx, strings.NewReader(in))
				assert.NoError(t, err)
				assert.Equal(t, 7, processed)
				assert.Equal(t, 4, accepted)
				assert.Equal(t, 3, discarded)
			},
		},
		{
			scenario: "Import data failed, malformed csv file",
			persisterFunc: func(g *domain.Geolocation) error {
				panic("should not be called")
			},
			assert: func(t *testing.T, uc domain.ImportGeolocationFromCSVFileUseCase) {
				in := `ip_address,country_code,country,city,latitude,longitude
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115`

				_, _, _, err := uc.Do(context.TODO(), strings.NewReader(in))
				assert.Error(t, err)
			},
		},
		{
			scenario: "Skipped all import data, persist database failed",
			persisterFunc: func(g *domain.Geolocation) error {
				return errors.New("DB error")
			},
			assert: func(t *testing.T, uc domain.ImportGeolocationFromCSVFileUseCase) {
				in := `ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115`

				processed, accepted, discarded, err := uc.Do(context.TODO(), strings.NewReader(in))
				assert.NoError(t, err)
				assert.Equal(t, 2, processed)
				assert.Equal(t, 0, accepted)
				assert.Equal(t, 2, discarded)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint
		t.Run(tc.scenario, func(t *testing.T) {
			pMock := domain.NewCallbackPersisterMock(tc.persisterFunc)

			uc := domain.NewImportGeolocationFromCSVFileToDBUseCase(pMock)

			tc.assert(t, uc)
		})
	}
}
