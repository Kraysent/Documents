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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

	g, _ := errgroup.WithContext(command.Context())

	g.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())
		return http.ListenAndServe(":2112", nil)
	})
	g.Go(func() error {
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

		return http.ListenAndServe(
			fmt.Sprintf("0.0.0.0:%d", command.Repository.Config.Server.Port), router,
		)
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
		log.Fatal("runtime error", zap.Error(err))
	}
}
