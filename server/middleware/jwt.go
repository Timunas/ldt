package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/timunas/ldt/server/handler"
	"github.com/timunas/ldt/server/token"
	"github.com/urfave/negroni"
)

func JwtMiddleware(tokenService token.Service) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		token, err := r.Cookie("Token")
		if err != nil {
			log.Info().Msgf("Cookie doesn't contain a valid authorization token")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := tokenService.ParseToken(token.Value)
		if err != nil {
			log.Info().Err(err).Msgf("Invalid token received")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(strings.TrimSpace(claims.Subject)) == 0 {
			log.Info().Msgf("No user id present in token")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), handler.RequestContextUserIDKey{}, claims.Subject)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
