package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/j0shgrant/leyline/internal"
	"github.com/oklog/run"
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

	// Read in env vars
	httpPort := internal.EnvInt("LEYLINE_FOREMAN_HTTP_PORT", 8081)
	region := internal.EnvString("LEYLINE_FOREMAN_REGION", "eu-west-1")
	apiKey := internal.EnvString("ABLY_API_KEY", "")
	if apiKey == "" {
		zap.S().Fatal("No environment variable set for ABLY_API_KEY")
	}

	// Concurrently run core flows
	g := run.Group{}

	// Configure HTTP routing
	g.Add(func() error {
		r := mux.NewRouter()
		r.HandleFunc("/ready", func(_ http.ResponseWriter, _ *http.Request) {})

		listenAddr := fmt.Sprintf(":%d", httpPort)

		zap.S().Infof("Listening on: [http://localhost:%s].", listenAddr)

		return http.ListenAndServe(listenAddr, r)
	}, func(err error) {
		zap.S().Fatal(err)
	})

	// Listen for Ably messages
	g.Add(func() error {
		channel, err := internal.NewAblyChannel(apiKey)
		if err != nil {
			return err
		}
		defer channel.Close()

		sub, err := channel.Subscribe("test")
		if err != nil {
			return err
		}
		defer sub.Close()

		zap.S().Infof("Listening for messages with region: [%s].", region)

		for serialisedMsg := range sub.MessageChannel() {
			var msg internal.Message
			err = msgpack.Unmarshal([]byte(serialisedMsg.Data.(string)), &msg)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			for _, msgRegion := range msg.Regions {
				if msgRegion == region {
					zap.S().Infof("Matched on region: [%s].", region)
					break
				}
			}
		}

		return nil
	}, func(err error) {
		zap.S().Fatal(err)
	})
	zap.S().Fatal(g.Run())
}


