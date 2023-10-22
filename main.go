package main

import (
	"fmt"
	"log"
	"net/http"

	"documents/internal/commands"
	"documents/internal/server"
	"documents/internal/server/handlers/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func main() {
	var command commands.Command
	if err := command.Init(); err != nil {
		log.Fatal("error during initialization", zap.Error(err))
	}
	defer func() {
		if err := command.Cleanup(); err != nil {
			log.Fatal("runtime error", zap.Error(err))
		}
	}()

	done := make(chan error)

	go func() {
		router := chi.NewRouter()
		router.Use(middleware.RequestID)
		router.Use(middleware.Recoverer)
		router.Use(middleware.Logger)
		router.Use(middleware.CleanPath)
		router.Use(command.Repository.SessionManager.LoadAndSave)

		for _, handler := range server.GetHandlers() {
			router.MethodFunc(handler.Method, handler.Path, handler.GetHandler(command.Repository))
		}

		router.Get("/auth/google/login", auth.GetGoogleLoginHandler(command.Repository))
		router.Get("/auth/google/callback", auth.GetGoogleCallbackHandler(command.Repository))

		if err := http.ListenAndServe(
			fmt.Sprintf("0.0.0.0:%d", command.Repository.Config.Server.Port), router,
		); err != nil {
			done <- err
		}
	}()

	if err := <-done; err != nil {
		log.Fatal("runtime error", zap.Error(err))
	}
}
