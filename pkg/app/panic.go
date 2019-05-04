package app

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dohernandez/geolocation-service/pkg/http/rest"
	requestid "github.com/dohernandez/geolocation-service/pkg/http/rest/request"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// PanicCatcher is a function to call when panic happens
type PanicCatcher func(rw http.ResponseWriter, req *http.Request, recover interface{}, stack []byte)

// PanicRecoverer recovers from panics and passes panic value to catchers
func PanicRecoverer(panicCatcher []PanicCatcher) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					stack := debug.Stack()
					for _, c := range panicCatcher {
						c(rw, req, rvr, stack)
					}

					return
				}
			}()

			next.ServeHTTP(rw, req)
		})
	}
}

// PanicPrinter is a PanicCatcher to print panic to STDERR, suitable for dev environment
func PanicPrinter(_ http.ResponseWriter, _ *http.Request, recover interface{}, stack []byte) {
	println("Panic:", fmt.Sprintf("%+v", recover))
	println(string(stack))
}

// PanicLogger is a PanicCatcher to log panic in http handling, suitable for staging and production environment
func PanicLogger(logger logrus.FieldLogger) PanicCatcher {
	return func(_ http.ResponseWriter, req *http.Request, recover interface{}, stack []byte) {
		logger.WithFields(logrus.Fields{
			"user-agent": req.UserAgent(),
			"path":       req.URL.Path,
			"remote":     req.RemoteAddr,
			"method":     req.Method,
			"uri":        req.URL.String(),
			"stack":      string(stack),
			"panic":      fmt.Sprintf("%+v", recover),
			"request-id": requestid.FromContext(req.Context()),
		}).Error("request panicked")
	}
}

// PanicResponse returns HTTP 500 (Internal Server Error) and panic message.
// Rendering errors (if any) will be logged.
// Consider hiding panic message if API is exposed to the public and may reveal sensitive info in panic message.
func PanicResponse(logger logrus.FieldLogger, hidePanicMessage bool) PanicCatcher {
	return func(rw http.ResponseWriter, req *http.Request, recover interface{}, stack []byte) {
		msg, ok := recover.(string)
		if !ok || hidePanicMessage {
			msg = "panic"
		}

		httpErr := &rest.ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      msg,
		}

		err := render.Render(rw, req, httpErr)
		if err != nil && logger != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("rendering error response")
		}
	}
}
