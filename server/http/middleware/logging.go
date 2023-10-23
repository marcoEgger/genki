package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/marcoEgger/genki/logger"
)

func LoggingHandler(handler http.Handler, skipEndpoint ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{ResponseWriter: w, body: bytes.NewBufferString("")}

		for _, skip := range skipEndpoint {
			if r.URL.String() == skip {
				handler.ServeHTTP(&sw, r)
				return
			}
		}

		log := logger.WithMetadata(r.Context()).WithFields(logger.Fields{
			"req.method": r.Method,
			"req.url":    r.URL,
		})

		// Read body in case we want to log it on request error.
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Warnf("failed to read body: %s", err)
		}
		// Reset body so the actual business logic can work with it
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Infof("incoming request to %s %s", r.Method, r.URL)
		defer func(started time.Time) {
			log = log.WithFields(logger.Fields{
				"took":   time.Since(started),
				"status": sw.status,
			})
			if sw.status == 400 || sw.status == 500 {
				log = log.WithFields(logger.Fields{
					"req.body": string(bodyBytes),
				})
				log.Warnf("served request to %s with error: %s", r.URL, sw.body.String())
			} else {
				log.Infof("served request to %s", r.URL)
			}
		}(time.Now())
		handler.ServeHTTP(&sw, r)
	})
}
