package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"              // для работы с table-сервисом
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"       // для работы с table-сервисом
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named" // для работы с table-сервисом
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"        // для работы с типами YDB и значениями
	yc "github.com/ydb-platform/ydb-go-yc"                     // для работы с YDB в Яндекс Облаке
)

//goland:noinspection GoUnusedExportedFunction
func Handler(rw http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()

	// Get YDB connection details from environment variables
	ydbEndpoint := os.Getenv("YDB_ENDPOINT")
	ydbDatabase := os.Getenv("YDB_DATABASE")

	if ydbEndpoint == "" || ydbDatabase == "" {
		log.Printf("Error: YDB_ENDPOINT or YDB_DATABASE environment variables not set")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "YDB configuration not found"}`))
		return
	}

	// Construct DSN from environment variables
	dsn := ydbEndpoint + "?database=" + ydbDatabase

	// создаем объект подключения db, является входной точкой для сервисов YDB
	db, err := ydb.Open(ctx, dsn,
		yc.WithMetadataCredentials(), // аутентификация изнутри виртуальной машины в Яндекс Облаке или из Яндекс Функции
	)
	if err != nil {
		log.Printf("Error connecting to YDB: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "Failed to connect to YDB"}`))
		return
	}
	// закрытие драйвера по окончании работы программы обязательно
	defer db.Close(ctx)

	var (
		readTx = table.TxControl(
			table.BeginTx(
				table.WithOnlineReadOnly(),
			),
			table.CommitTx(),
		)
	)
	var user struct {
		ID   int32   `json:"id"`
		Name *string `json:"name"` // optional
	}
	err = db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			var (
				res  result.Result
				id   int32   // переменная для required результатов
				name *string // указатель - для опциональных результатов
			)
			_, res, err = s.Execute(
				ctx,
				readTx,
				`
        DECLARE $id AS Int32;
        SELECT
          id,
          name,
        FROM
          users
        WHERE
          id = $id;
      `,
				table.NewQueryParameters(
					table.ValueParam("$id", types.Int32Value(3)), // подстановка в условие запроса
				),
			)
			if err != nil {
				return err
			}
			defer res.Close() // закрытие result'а обязательно
			log.Printf("> select_simple_transaction:\n")
			for res.NextResultSet(ctx) {
				for res.NextRow() {
					// в ScanNamed передаем имена колонок из строки сканирования,
					// адреса (и типы данных), куда следует присвоить результаты запроса
					err = res.ScanNamed(
						named.Required("id", &id),
						named.Optional("name", &name),
					)
					if err != nil {
						return err
					}
					log.Printf(
						"  > %d %s %s\n",
						id, *name,
					)
					user.ID = id
					user.Name = name
				}
			}
			return res.Err()
		},
	)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "Failed to execute query"}`))
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "Failed to marshal response"}`))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	_, err = rw.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}
