package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/timunas/ldt/server/model"
)

func TodosHandler(repository model.TodoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			userID := r.Context().Value(RequestContextUserIDKey{}).(string)
			result, err := repository.FindByUserID(userID)

			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch all objects...")
				w.WriteHeader(http.StatusInternalServerError)
			}

			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Error().Err(err).Msg("Failed to encode response body...")
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "POST":
			userID := r.Context().Value(RequestContextUserIDKey{}).(string)
			body := TodoRequest{}
			err := json.NewDecoder(r.Body).Decode(&body)

			if err != nil {
				log.Error().Err(err).Msg("Failed to decode request body...")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			result, err := repository.Save(model.NewTodo(body.Name, body.Description, userID))

			if err != nil {
				log.Error().Err(err).Msg("Failed to save object...")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Error().Err(err).Msg("Failed to encode response body...")
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Header().Add("location", r.RequestURI+"/"+result.ID)
				w.WriteHeader(http.StatusCreated)
			}
		}
	}
}
