package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/matsuboshi/league-matrix-app/internal/handler"
)

const port = "8080"

func main() {
	matrixHandler := handler.NewMatrixHandler()

	http.HandleFunc("/", matrixHandler.ListMatrixOperations)
	http.HandleFunc("/matrix", matrixHandler.ListMatrixOperations)
	http.HandleFunc("/matrix/", matrixHandler.ProcessMatrix)

	slog.Info("starting HTTP server", "port", port, "address", "http://localhost:"+port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("server failed to start", "error", err, "port", port)
		os.Exit(1)
	}
}
