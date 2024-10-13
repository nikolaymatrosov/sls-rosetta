package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/tracing"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type iamToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

type iamRequestMiddleware struct {
	options sesv2.Options
}

func (*iamRequestMiddleware) ID() string {
	return "IamToken"
}

func (m *iamRequestMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	_, span := tracing.StartSpan(ctx, "IamToken")
	defer span.End()

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unexpected transport type %T", in.Request)
	}

	tokenValue := ctx.Value("lambdaRuntimeTokenJSON")
	tokenStr, ok := tokenValue.(string)
	if !ok {
		return out, metadata, fmt.Errorf("unexpected token type %T", tokenValue)
	}

	var token iamToken
	err = json.Unmarshal([]byte(tokenStr), &token)
	if err != nil {
		return middleware.FinalizeOutput{}, middleware.Metadata{}, err
	}

	req.Header.Set("X-YaCloud-SubjectToken", token.Token)

	span.End()
	return next.HandleFinalize(ctx, in)
}

func swapAuth() func(options *sesv2.Options) {
	return func(options *sesv2.Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			_, err := stack.Finalize.Swap("Signing", &iamRequestMiddleware{})
			if err != nil {
				return err
			}
			return nil
		})
	}
}
