package main

import (
	"context"
	"log"

	websocketapi "github.com/yandex-cloud/go-genproto/yandex/cloud/serverless/apigateway/websocket/v1"
	websocketsdk "github.com/yandex-cloud/go-sdk/services/serverless/apigateway/websocket/v1"
	"github.com/yandex-cloud/go-sdk/v2"
	"github.com/yandex-cloud/go-sdk/v2/credentials"
	"github.com/yandex-cloud/go-sdk/v2/pkg/options"
)

func Handler() {
	ctx := context.Background()

	sdk, err := ycsdk.Build(ctx,
		options.WithCredentials(credentials.InstanceServiceAccount()),
	)
	if err != nil {
		log.Fatal(err)
	}
	cc := websocketsdk.NewConnectionClient(sdk)
	connection, err := cc.Get(ctx, &websocketapi.GetConnectionRequest{
		ConnectionId: "connection-id",
	})
	if err != nil {
		return
	}

	log.Printf("Connection: %v", connection)

}
