package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mobamoh/twitter-go-graphql/config"
	"github.com/mobamoh/twitter-go-graphql/domain"
	"github.com/mobamoh/twitter-go-graphql/graph"
	"github.com/mobamoh/twitter-go-graphql/jwt"
	"github.com/mobamoh/twitter-go-graphql/postgres"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	config.LoadEnv(".env")
	conf := config.New()

	db := postgres.New(ctx, conf)
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Timeout(time.Second * 60))

	userRepo := postgres.NewUserRepo(db)
	authTokenService := jwt.NewTokenService(conf)
	authService := domain.NewAuthService(userRepo, authTokenService)

	tweetRepo := postgres.NewTweetRepo(db)
	tweetService := domain.NewTweetService(tweetRepo)

	router.Use(graph.DataloaderMiddleware(&graph.Repos{
		UserRepo: userRepo,
	}))
	router.Use(authMiddleware(authTokenService))
	router.Handle("/", playground.AltairHandler("Go Twitter Clone", "/query"))
	router.Handle("/query", handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					AuthService:  authService,
					TweetService: tweetService,
				},
			},
		),
	))

	log.Fatal(http.ListenAndServe(":8080", router))
}
