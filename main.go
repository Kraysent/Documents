package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"documents/internal/actions"
	"documents/internal/commands"
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
		http.HandleFunc("/api/v1/document/insert", func(writer http.ResponseWriter, request *http.Request) {
			if err := actions.InsertDocument(ctx, command.Repository); err != nil {
				writer.Write([]byte(fmt.Sprintf(err.Error())))
			}
		})
		if err := http.ListenAndServe(
			fmt.Sprintf(":%d", command.Repository.Config.Server.Port), nil,
		); err != nil {
			done <- err
		}
	}()

	if err := <-done; err != nil {
		log.Fatal("runtime error", zap.Error(err))
	}

}
