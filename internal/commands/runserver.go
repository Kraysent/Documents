package commands

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/log"
	"documents/internal/server"
	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	ConfigEnv = "CONFIG"
)

type Command struct {
	ctx           context.Context
	Repository    *core.Repository
	APIServer     *http.Server
	MetricsServer *http.Server
}

func (c *Command) Context() context.Context {
	if c.ctx == nil {
		c.ctx = context.Background()
	}

	return c.ctx
}

func (c *Command) initAPIServer() error {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)
	router.Use(web.CORSMiddleware(c.Repository))
	if c.Repository.Config.Metrics.Name == "" {
		// This is needed in order to avoid duplication of metrics registration in parallel tests
		middlewareName, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		c.Repository.Config.Metrics.Name = middlewareName.String()
	}

	router.Use(chiprometheus.NewMiddleware(c.Repository.Config.Metrics.Name))
	router.Use(c.Repository.SessionManager.LoadAndSave)

	for _, handler := range server.GetHandlers() {
		router.MethodFunc(handler.Method, handler.Path, handler.GetHandler(c.Repository))
	}

	for _, handler := range server.GetAuthHandlers() {
		router.MethodFunc(handler.Method, handler.Path, handler.GetHandler(c.Repository))
	}

	var srv http.Server
	srv.Addr = fmt.Sprintf("0.0.0.0:%d", c.Repository.Config.Server.Port)
	srv.Handler = router

	c.APIServer = &srv

	return nil
}

func (c *Command) initMetricsServer() {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	var srv http.Server
	srv.Addr = ":2112"
	srv.Handler = router

	c.MetricsServer = &srv
}

func (c *Command) Init() error {
	configPath, ok := os.LookupEnv(ConfigEnv)
	if !ok {
		return fmt.Errorf("no config specified")
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s env set but file %s does not exist", ConfigEnv, configPath)
	}

	config, err := core.ParseConfig(configPath)
	if err != nil {
		return err
	}

	if err := log.InitLogger(); err != nil {
		return err
	}

	repo, err := core.NewRepository(config)
	if err != nil {
		return err
	}

	c.Repository = repo

	if err := c.Repository.Storage.Connect(c.Context()); err != nil {
		return err
	}

	c.Repository.SessionManager.Store = pgxstore.New(c.Repository.Storage.GetDB())

	if err := c.initAPIServer(); err != nil {
		return err
	}
	c.initMetricsServer()

	log.Info("Connected to the database")
	return nil
}

func (c *Command) Start() error {
	g, _ := errgroup.WithContext(c.Context())

	g.Go(func() error {
		return c.MetricsServer.ListenAndServe()
	})
	g.Go(func() error {
		log.Info("Starting server",
			zap.Int("port", c.Repository.Config.Server.Port),
			zap.String("url", fmt.Sprintf("http://%s:%d",
				c.Repository.Config.Server.Host, c.Repository.Config.Server.Port)))

		return c.APIServer.ListenAndServe()
	})

	return g.Wait()
}

func (c *Command) Cleanup() error {
	if err := c.MetricsServer.Shutdown(c.Context()); err != nil {
		return err
	}
	log.Info("Metrics server shut down")
	if err := c.APIServer.Shutdown(c.Context()); err != nil {
		return err
	}
	log.Info("API server shut down")
	if err := c.Repository.Storage.Disconnect(c.Context()); err != nil {
		return err
	}
	log.Info("Disconnected from database")

	return nil
}
