package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"github.com/timunas/ldt/server/google"
	"github.com/timunas/ldt/server/model"
	"github.com/timunas/ldt/server/token"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func AuthHandler(tokenService token.Service, googleConfig *oauth2.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		claims := jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "ldt",
		}
		state, err := tokenService.NewToken(claims)
		if err != nil {
			log.Error().Err(err).Msgf("Failed while generating state token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authCodeURL := googleConfig.AuthCodeURL(state)
		http.Redirect(w, r, authCodeURL, http.StatusTemporaryRedirect)
	}
}

func AuthCallbackHandler(repository model.UserRepository, tokenService token.Service, googleConfig *oauth2.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_, err := tokenService.ParseToken(r.URL.Query().Get("state"))

		if err != nil {
			log.Error().Err(err).Msgf("Invalid state received")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		code := r.URL.Query().Get("code")
		googleToken, err := googleConfig.Exchange(context.Background(), code)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to exchange code")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		client := googleConfig.Client(context.Background(), googleToken)
		response, err := client.Get(google.GoogleUserInfoEndpoint)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to retrieve user information")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		defer response.Body.Close()
		userInfo := UserInfo{}

		err = json.NewDecoder(response.Body).Decode(&userInfo)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse user information response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := repository.FindByEmail(userInfo.Email)
		if err != nil {
			log.Info().Err(err).Msg("User not found. Creating a new one...")
			user, err = repository.Save(model.NewUser(userInfo.Name, userInfo.Email))
			if err != nil {
				log.Error().Err(err).Msg("Failed to save object...")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		expiration := time.Now().Add(time.Hour * 1)
		claims := jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiration.Unix(),
			Issuer:    "ldt",
			Subject:   user.ID,
		}
		token, err := tokenService.NewToken(claims)
		if err != nil {
			log.Error().Err(err).Msgf("Failed while generating cookie token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "Token",
			Value:    token,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Secure:   true,
			Expires:  expiration,
		}
		http.SetCookie(w, &cookie)

		err = json.NewEncoder(w).Encode(userInfo)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encode response body...")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
