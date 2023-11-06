package main

import (
	"fmt"
	"net/http"

	"documents/internal/commands"
	"documents/internal/library/web"
	"documents/internal/log"
	"documents/internal/server"
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
		router.Use(web.CORSMiddleware)
		router.Use(command.Repository.SessionManager.LoadAndSave)

		for _, handler := range server.GetHandlers() {
			router.MethodFunc(handler.Method, handler.Path, handler.GetHandler(command.Repository))
		}

		for _, handler := range server.GetAuthHandlers() {
			router.MethodFunc(handler.Method, handler.Path, handler.GetHandler(command.Repository))
		}

		log.Info("Starting server",
			zap.Int("port", command.Repository.Config.Server.Port),
			zap.String("url", fmt.Sprintf("http://%s:%d",
				command.Repository.Config.Server.Host, command.Repository.Config.Server.Port)))

		if err := http.ListenAndServe(
			fmt.Sprintf("0.0.0.0:%d", command.Repository.Config.Server.Port), router,
		); err != nil {
			log.Info("Stopping server", zap.Int("port", command.Repository.Config.Server.Port))
			done <- err
		}
	}()

	if err := <-done; err != nil {
		fmt.Println(err)
		log.Fatal("runtime error", zap.Error(err))
	}
}
