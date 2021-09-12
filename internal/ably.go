package internal

import (
	"github.com/ably/ably-go/ably"
)

const channelName = "test"

func NewAblyChannel(apiKey string) (*ably.RealtimeChannel, error){
	client, err := ably.NewRealtimeClient(ably.NewClientOptions(apiKey))
	if err != nil {
		return nil, err
	}

	return client.Channels.Get(channelName), nil
}
