package main

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/smithy-go/endpoints"
)

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params sesv2.EndpointParameters) (
	transport.Endpoint, error,
) {
	u, err := url.Parse("https://postbox.cloud.yandex.net")
	if err != nil {
		return transport.Endpoint{}, err
	}
	return transport.Endpoint{
		URI: *u,
	}, nil
}
