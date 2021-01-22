package app

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qiniu/qmgo"
	"github.com/rs/zerolog/log"
	"github.com/timunas/ldt/server/handler"
	"github.com/timunas/ldt/server/middleware"
	"github.com/timunas/ldt/server/repository"
	"github.com/urfave/negroni"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})

	if err != nil {
		log.Panic().Err(err).Msg("Failed connecting to database")
	}

	defer func() {
		if err = client.Close(ctx); err != nil {
			log.Panic().Err(err).Msg("Failed while disconnecting from database")
		}
	}()

	db := client.Database("ldt")
	todoCollection := db.Collection("todo")

	r := mux.NewRouter()
	r.HandleFunc("/todos", handler.TodosHandler(repository.NewTodoRepository(todoCollection, &ctx))).Methods("GET", "POST")
	r.HandleFunc("/todos/{id}", handler.TodoHandler(repository.NewTodoRepository(todoCollection, &ctx))).Methods("GET", "DELETE", "PUT")

	n := negroni.New().With(
		negroni.NewRecovery(),
		negroni.HandlerFunc(middleware.LoggingMiddleware),
	)
	n.UseHandler(r)

	port := 8080
	log.Info().Msgf("Starting server at port: %d", port)
	log.Error().Err(http.ListenAndServe(":"+strconv.Itoa(port), n))
}
