package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"documents/internal/commands"
	"documents/internal/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	var command commands.Command
	if err := command.Init(); err != nil {
		log.Fatal("error during initialization", zap.Error(err))
	}

	done := make(chan error)

	go func() {
		if err := command.Repository.Storage.DocumentStorage.Connect(ctx); err != nil {
			done <- err
		}

		router := chi.NewRouter()
		router.Use(middleware.RequestID)
		router.Use(middleware.Recoverer)
		router.Use(middleware.Logger)
		router.Use(middleware.CleanPath)

		for _, handler := range server.GetHandlers() {
			router.MethodFunc(handler.Method, handler.Path, handler.GetHandler(command.Repository))
		}

		if err := http.ListenAndServe(
			fmt.Sprintf("0.0.0.0:%d", command.Repository.Config.Server.Port), router,
		); err != nil {
			done <- err
		}

		if err := command.Repository.Storage.DocumentStorage.Disconnect(ctx); err != nil {
			done <- err
		}
	}()

	if err := <-done; err != nil {
		log.Fatal("runtime error", zap.Error(err))
	}
}
