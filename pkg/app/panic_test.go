package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dohernandez/geolocation-service/pkg/app"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestPanicRecoverer_Middleware(t *testing.T) {
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		panic("AAA!")
	})

	gotPanic := false

	panicatchers := []app.PanicCatcher{
		func(rw http.ResponseWriter, req *http.Request, recover interface{}, stack []byte) {
			gotPanic = true
			assert.Equal(t, interface{}("AAA!"), recover)
		},
	}

	r := chi.NewRouter()
	r.Use(app.PanicRecoverer(panicatchers))
	r.Method(http.MethodGet, "/panic", h)

	r.ServeHTTP(nil, httptest.NewRequest("", "/panic", nil))
	assert.True(t, gotPanic)
}
