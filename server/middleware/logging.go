package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func LoggingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Info().
		Str("method", r.Method).
		Str("path", r.RequestURI).
		Msg("Request received...")

	next.ServeHTTP(rw, r)

	statusCode := rw.(negroni.ResponseWriter).Status()
	log.Info().
		Str("method", r.Method).
		Str("path", r.RequestURI).
		Int("status_code", statusCode).
		Msgf("Response sent: %s", http.StatusText(statusCode))
}
