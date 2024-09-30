package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

func Long(ctx context.Context, request []byte) (*Response, error) {
	inFlightRequests++
	handelesdRequests++
	// Создание логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	time.Sleep(1 * time.Second)

	requestBody := &RequestBody{}
	// Массив байтов, содержащий тело запроса, преобразуется в соответствующий объект
	err := json.Unmarshal(request, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing request: %v", err)
	}
	headers := make(map[string]string)
	for k, v := range requestBody.Headers {
		headers[strings.ToLower(k)] = v
	}

	requestId, ok := headers["x-request-id"]
	if !ok {
		return nil, fmt.Errorf("requestId not found")
	}

	// В журнале будет напечатано название HTTP-метода, с помощью которого осуществлен запрос, а также тело запроса
	logger.Info("got request",
		zap.String("requestId", requestId),
		zap.Int("inFlightRequests", inFlightRequests),
		zap.Int("concurrency", concurrency),
		zap.Int("handledRequests", handelesdRequests),
	)

	req := &Request{}
	// Поле body запроса преобразуется в объект типа Request для получения переданного имени
	err = json.Unmarshal([]byte(requestBody.Body), &req)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing body: %v", err)
	}

	name := req.Name

	inFlightRequests--
	// Возвращается объект типа Response, содержащий код состояния 200 и приветствие с именем
	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("Hello, %s", name),
	}, nil
}
