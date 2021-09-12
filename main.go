package main

import (
	"fmt"
	"github.com/ably/ably-go/ably"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to create logger.")
		os.Exit(-1)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	client, err := ably.NewRealtimeClient(ably.NewClientOptions("lBQ6LQ.o1BZBQ:T-J49Y28bHvLTx9Q"))
	if err != nil {
		zap.S().Fatal(err)
	}

	client.Channels.Get("test")

	for _, channel := range client.Channels.All() {
		zap.S().Info(channel.Name)
	}

	_ = client.Channels.Release("test")
}
