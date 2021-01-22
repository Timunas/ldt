package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/timunas/ldt/server/model"
)

func TodoHandler(repository model.TodoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			id := mux.Vars(r)["id"]

			result, err := repository.FindByID(id)

			if err != nil {
				log.Error().Err(err).Msgf("Tried to fetch object with id: %s", id)
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w)
				return
			}

			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Error().Err(err).Msg("Failed to encode response body...")
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case "DELETE":
			id := mux.Vars(r)["id"]

			result, err := repository.FindByID(id)

			if err != nil {
				log.Error().Err(err).Msgf("Tried to fetch object with id: %s", id)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			err = repository.Delete(result)

			if err != nil {
				log.Error().Err(err).Msgf("Failed deleting object with id: %s", id)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		case "PUT":
			id := mux.Vars(r)["id"]

			todo, err := repository.FindByID(id)

			if err != nil {
				log.Error().Err(err).Msgf("Tried to fetch object with id: %s", id)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			body := TodoRequest{}
			err = json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				log.Error().Err(err).Msg("Failed to decode request body...")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			todo.Name = body.Name
			todo.Description = body.Description

			result, err := repository.Save(todo)

			if err != nil {
				log.Error().Err(err).Msgf("Failed updating object with id: %s", id)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Error().Err(err).Msg("Failed to encode response body...")
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}
