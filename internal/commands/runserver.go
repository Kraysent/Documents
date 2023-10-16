package commands

import (
	"context"
	"fmt"
	"os"

	"documents/internal/core"
	"documents/internal/log"
	"documents/internal/server"
)

type Command struct {
	ctx        context.Context
	Repository *core.Repository
	Server     *server.Server
}

func (c *Command) Context() context.Context {
	if c.ctx == nil {
		c.ctx = context.Background()
	}

	return c.ctx
}

func (c *Command) Init() error {
	configPath, ok := os.LookupEnv("CONFIG")
	if !ok {
		return fmt.Errorf("no config specified")
	}

	config, err := core.ParseConfig(configPath)
	if err != nil {
		return err
	}

	repo, err := core.NewRepository(config)
	if err != nil {
		return err
	}

	c.Repository = repo

	c.Server = server.NewServer(repo)

	return log.InitLogger(config.Logging.StdoutPath, config.Logging.StderrPath)
}
