package main

import (
	"fmt"

	"documents/internal/commands"
	"documents/internal/log"
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

	if err := command.Start(); err != nil {
		fmt.Println(err)
		log.Fatal("runtime error", zap.Error(err))
	}
}
