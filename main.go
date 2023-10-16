package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
		if err := command.Repository.Storage.DocumentStorage.Connect(ctx); err != nil {
			done <- err
		}

		http.HandleFunc("/api/v1/document/insert", command.Server.InsertDocument)
		http.HandleFunc("/api/v1/document/get", command.Server.GetDocumentByID)
		http.HandleFunc("/api/v1/document/get/username", command.Server.GetDocumentByUsernameAndType)
		if err := http.ListenAndServe(
			fmt.Sprintf(":%d", command.Repository.Config.Server.Port), nil,
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
