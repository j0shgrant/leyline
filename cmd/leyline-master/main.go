package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/j0shgrant/leyline/internal"
	"github.com/vmihailenco/msgpack/v5"
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

	// load apiKey
	apiKey := internal.EnvString("ABLY_API_KEY", "")
	if apiKey == "" {
		zap.S().Fatal("No environment variable set for ABLY_API_KEY")
	}

	// Configure HTTP routing
	r := mux.NewRouter()
	ablyHandler, err := newAblyMasterHandler(apiKey)
	r.HandleFunc("/ready", func(_ http.ResponseWriter, _ *http.Request) {})
	r.HandleFunc("/test/{region}", ablyHandler)

	listenAddr := fmt.Sprintf(":%d", 8080)
	zap.S().Infof("Listening on: [http://localhost:%s].", listenAddr)
	zap.S().Fatal(http.ListenAndServe(listenAddr, r))
}

func newAblyMasterHandler(apiKey string) (func(w http.ResponseWriter, r *http.Request), error) {
	channel, err := internal.NewAblyChannel(apiKey)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if region, exists := vars["region"]; exists {
			zap.S().Info(region)

			msg := internal.Message{
				Regions: []string{region},
			}

			serialisedMsg, err := msgpack.Marshal(msg)
			if err != nil {
				zap.S().Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))

				return
			}

			_, err = channel.Publish("test", serialisedMsg)
			if err != nil {
				zap.S().Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))

				return
			}

			zap.S().Info("Pushed a message to Ably")

			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("must pass in a region"))

		return
	}, nil
}
