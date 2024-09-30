package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	environ "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"go.uber.org/zap"
)

var db *ydb.Driver

func YdbHandler(ctx context.Context, request []byte) (*Response, error) {
	inFlightRequests++
	handelesdRequests++
	// Создание логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	var err error
	dsn := os.Getenv("YDB_DSN")
	if db == nil {
		db, err = ydb.Open(ctx, dsn,
			environ.WithEnvironCredentials(),
			ydb.WithDialTimeout(time.Second),
		)
		if err != nil {
			return nil, err
		}
	}

	requestBody := &RequestBody{}
	// Массив байтов, содержащий тело запроса, преобразуется в соответствующий объект
	err = json.Unmarshal(request, &requestBody)
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

	query := `
		DECLARE $id AS Text;
		DECLARE	$data AS JSONDocument;

		INSERT INTO requests ( id, data ) VALUES ( $id, $data );
	`

	req := &Request{}
	// Поле body запроса преобразуется в объект типа Request для получения переданного имени
	err = json.Unmarshal([]byte(requestBody.Body), &req)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing body: %v", err)
	}

	err = db.Table().Do(ctx,
		func(ctx context.Context, session table.Session) (err error) {
			_, _, err = session.Execute(
				ctx,
				table.SerializableReadWriteTxControl(table.CommitTx()),
				query,
				table.NewQueryParameters(
					table.ValueParam("$id", types.TextValue(requestId)),
					table.ValueParam("$data", types.JSONDocumentValue(requestBody.Body)),
				),
			)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	// В логе будет напечатано название HTTP-метода, с помощью которого осуществлен запрос, а также тело запроса
	logger.Info("got request",
		zap.String("requestId", requestId),
		zap.Int("inFlightRequests", inFlightRequests),
		zap.Int("concurrency", concurrency),
		zap.Int("handledRequests", handelesdRequests),
	)

	name := req.Name

	inFlightRequests--
	// Возвращается объект типа Response, содержащий код состояния 200 и приветствие с именем
	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("Hello, %s", name),
	}, nil
}
