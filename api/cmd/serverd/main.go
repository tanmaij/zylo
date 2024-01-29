package main

import (
	"log"
	"net/http"

	"github.com/tanmaij/zylo/cmd/serverd/router"
	"github.com/tanmaij/zylo/config"
	"github.com/tanmaij/zylo/pkg/redis"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := config.Initial(); err != nil {
		log.Fatal("", err)
	}

	appCfg := config.Instance.App

	_, err := redis.NewRedisClient()
	if err != nil {
		log.Fatal("")
	}

	r := chi.NewRouter()
	appRouter := router.Router{}
	appRouter.Routes(r)

	http.ListenAndServe(appCfg.Port, r)
}
