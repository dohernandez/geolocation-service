package domain_test

import (
	"context"
	"strings"
	"testing"

	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/pkg/log"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestImportGeolocationFromCSVFileUseCase(t *testing.T) {
	testCases := []struct {
		scenario string

		assert func(t *testing.T, uc domain.ImportGeolocationFromCSVFileUseCase)
	}{
		{
			scenario: "Import data successfully",
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

				err := uc.Do(ctx, strings.NewReader(in))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint
		t.Run(tc.scenario, func(t *testing.T) {
			pMock := domain.NewCallbackPersisterMock(func(g *domain.Geolocation) error {
				// checking that the geolocation object is valid
				assert.NoError(t, g.Validate())
				// checking that the geolocation object has ID
				assert.NotEqual(t, uuid.Nil, g.ID)

				return nil
			})

			uc := domain.NewImportGeolocationFromCSVFileToDBUseCase(pMock)

			tc.assert(t, uc)
		})
	}
}
