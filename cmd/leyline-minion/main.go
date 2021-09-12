package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {
	// Configure logging
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to create logger.")
		os.Exit(-1)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Configure HTTP routing
	r := mux.NewRouter()
	r.HandleFunc("/ready", func(_ http.ResponseWriter, _ *http.Request) {})

	listenAddr := fmt.Sprintf(":%d", 8090)

	zap.S().Infof("Listening on: [http://localhost:%s].", listenAddr)
	zap.S().Fatal(http.ListenAndServe(listenAddr, r))
}
