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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	ConfigEnv = "CONFIG"
)

type Command struct {
	ctx        context.Context
	Repository *core.Repository
}

func (c *Command) Context() context.Context {
	if c.ctx == nil {
		c.ctx = context.Background()
	}

	return c.ctx
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

	log.Info("Connected to the database")
	return nil
}

func (c *Command) Start() error {
	g, _ := errgroup.WithContext(c.Context())

	g.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())
		return http.ListenAndServe(":2112", nil)
	})
	g.Go(func() error {
		router := chi.NewRouter()
		router.Use(middleware.RequestID)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(middleware.CleanPath)
		router.Use(web.CORSMiddleware(c.Repository))
		router.Use(chiprometheus.NewMiddleware("documents"))
		router.Use(c.Repository.SessionManager.LoadAndSave)

		for _, handler := range server.GetHandlers() {
			router.MethodFunc(handler.Method, handler.Path, handler.GetHandler(c.Repository))
		}

		for _, handler := range server.GetAuthHandlers() {
			router.MethodFunc(handler.Method, handler.Path, handler.GetHandler(c.Repository))
		}

		log.Info("Starting server",
			zap.Int("port", c.Repository.Config.Server.Port),
			zap.String("url", fmt.Sprintf("http://%s:%d",
				c.Repository.Config.Server.Host, c.Repository.Config.Server.Port)))

		return http.ListenAndServe(
			fmt.Sprintf("0.0.0.0:%d", c.Repository.Config.Server.Port), router,
		)
	})

	return g.Wait()
}

func (c *Command) Cleanup() error {
	if err := c.Repository.Storage.Disconnect(c.Context()); err != nil {
		return err
	}
	log.Info("Disconnected from the database")

	return nil
}
