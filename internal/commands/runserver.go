package commands

import (
	"context"
	"fmt"
	"os"

	"documents/internal/core"
	"documents/internal/log"
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
	configPath, ok := os.LookupEnv("CONFIG")
	if !ok {
		return fmt.Errorf("no config specified")
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

	log.Info("Connected to the database")
	return nil
}

func (c *Command) Cleanup() error {
	if err := c.Repository.Storage.Disconnect(c.Context()); err != nil {
		return err
	}
	log.Info("Disconnected from the database")

	return nil
}
