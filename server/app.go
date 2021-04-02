package app

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/rs/zerolog/log"
	"github.com/timunas/ldt/server/google"
	"github.com/timunas/ldt/server/handler"
	"github.com/timunas/ldt/server/middleware"
	"github.com/timunas/ldt/server/repository"
	"github.com/timunas/ldt/server/token"
	"github.com/urfave/negroni"
	"golang.org/x/oauth2"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	// Environment setup
	err := godotenv.Load()
	if err != nil {
		log.Panic().Err(err).Msg("Failed reading env file")
	}

	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")

	// Database setup
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

	userRepository := buildUserRepository(db, &ctx)
	todoRepository := buildTodoRepository(db, &ctx)

	// JWT service configuration
	jwtPrivateKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))
	jwtPublicKey := []byte(os.Getenv("JWT_PUBLIC_KEY"))
	tokenService, err := token.NewJwtService(jwtPrivateKey, jwtPublicKey)
	if err != nil {
		log.Panic().Err(err).Msg("Failed initializing JWT service")
	}

	// Google client configuration
	googleConfig := buildGoogleConfiguration(host)

	// Router configuration
	router := mux.NewRouter()
	router.HandleFunc("/auth", handler.AuthHandler(tokenService, googleConfig)).Methods("GET")
	router.HandleFunc("/auth/callback", handler.AuthCallbackHandler(userRepository, tokenService, googleConfig)).Methods("GET")
	router.Handle(
		"/todos",
		negroni.New(
			middleware.JwtMiddleware(tokenService),
			negroni.WrapFunc(handler.TodosHandler(todoRepository)),
		),
	).Methods("GET", "POST")
	router.Handle(
		"/todos/{id}",
		negroni.New(
			middleware.JwtMiddleware(tokenService),
			negroni.WrapFunc(handler.TodoHandler(todoRepository)),
		),
	).Methods("GET", "DELETE", "PUT")

	n := negroni.New().With(
		negroni.NewRecovery(),
		negroni.HandlerFunc(middleware.LoggingMiddleware),
	)
	n.UseHandler(router)

	log.Info().Msgf("Starting server at port: %s:%s", host, port)
	log.Error().Err(http.ListenAndServe(":"+port, n))
}

func buildGoogleConfiguration(host string) *oauth2.Config {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	var redirectHost string
	if os.Getenv("IS_DEV") == "true" {
		redirectHost = "http://localhost:8080"
	} else {
		redirectHost = "https://" + host
	}

	return google.NewConfig(googleClientID, googleClientSecret, redirectHost+"/auth/callback")
}

func buildUserRepository(db *qmgo.Database, ctx *context.Context) *repository.UserRepo {
	userCollection := db.Collection("user")
	err := userCollection.CreateOneIndex(context.Background(), options.IndexModel{Key: []string{"email"}, Unique: true})

	if err != nil {
		log.Panic().Err(err).Msg("Failed creating user collection index")
	}

	return repository.NewUserRepository(userCollection, ctx)
}

func buildTodoRepository(db *qmgo.Database, ctx *context.Context) *repository.TodoRepo {
	todoCollection := db.Collection("todo")
	err := todoCollection.CreateOneIndex(context.Background(), options.IndexModel{Key: []string{"user_id"}})
	if err != nil {
		log.Panic().Err(err).Msg("Failed creating todo collection index")
	}
	return repository.NewTodoRepository(todoCollection, ctx)
}
