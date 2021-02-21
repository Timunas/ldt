package app

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/qiniu/qmgo"
	"github.com/rs/zerolog/log"
	"github.com/timunas/ldt/server/google"
	"github.com/timunas/ldt/server/handler"
	"github.com/timunas/ldt/server/middleware"
	"github.com/timunas/ldt/server/repository"
	"github.com/timunas/ldt/server/token"
	"github.com/urfave/negroni"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	err := godotenv.Load()
	if err != nil {
		log.Panic().Err(err).Msg("Failed reading env file")
	}

	mongoURI := os.Getenv("MONGO_URI")

	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: mongoURI})

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

	jwtPrivateKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))
	jwtPublicKey := []byte(os.Getenv("JWT_PUBLIC_KEY"))
	tokenService, err := token.NewJwtService(jwtPrivateKey, jwtPublicKey)
	if err != nil {
		log.Panic().Err(err).Msg("Failed initializing JWT service")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleConfig := google.NewConfig(googleClientID, googleClientSecret, "http://localhost:8080/auth/callback")

	r := mux.NewRouter()
	r.HandleFunc("/auth", handler.AuthHandler(tokenService, googleConfig)).Methods("GET")
	r.HandleFunc("/auth/callback", handler.AuthCallbackHandler(tokenService, googleConfig)).Methods("GET")
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
