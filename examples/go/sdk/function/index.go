package main

import (
	"context"
	"log"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/serverless/apigateway/websocket/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

func Handler() {

	ctx := context.Background()

	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: ycsdk.InstanceServiceAccount(),
	})
	if err != nil {
		log.Fatal(err)
	}
	connection, err := sdk.Serverless().APIGatewayWebsocket().Connection().Get(ctx, &websocket.GetConnectionRequest{
		ConnectionId: "connection-id",
	})
	if err != nil {
		return
	}

	log.Printf("Connection: %v", connection)

}
