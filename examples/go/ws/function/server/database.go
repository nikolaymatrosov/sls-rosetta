package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

// StoreConnection stores or updates a connection in YDB
// Implements one-per-user policy by deleting existing connections for the user first
func StoreConnection(ctx context.Context, db *ydb.Driver, connectionID, userID string) error {
	return db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		// First, delete any existing connections for this user (one-per-user policy)
		deleteQuery := `
			DECLARE $user_id AS Utf8;
			DELETE FROM connections
			WHERE user_id = $user_id;
		`

		_, _, err := s.Execute(ctx, table.DefaultTxControl(), deleteQuery,
			table.NewQueryParameters(
				table.ValueParam("$user_id", types.UTF8Value(userID)),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to delete existing connection: %w", err)
		}

		// Now insert the new connection
		insertQuery := `
			DECLARE $connection_id AS Utf8;
			DECLARE $user_id AS Utf8;
			DECLARE $connected_at AS Timestamp;

			UPSERT INTO connections (connection_id, user_id, connected_at)
			VALUES ($connection_id, $user_id, $connected_at);
		`

		_, _, err = s.Execute(ctx, table.DefaultTxControl(), insertQuery,
			table.NewQueryParameters(
				table.ValueParam("$connection_id", types.UTF8Value(connectionID)),
				table.ValueParam("$user_id", types.UTF8Value(userID)),
				table.ValueParam("$connected_at", types.TimestampValueFromTime(time.Now())),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to insert connection: %w", err)
		}

		return nil
	})
}

// RemoveConnection removes a connection by user ID
func RemoveConnection(ctx context.Context, db *ydb.Driver, userID string) error {
	return db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := `
			DECLARE $user_id AS Utf8;
			DELETE FROM connections
			WHERE user_id = $user_id;
		`

		_, _, err := s.Execute(ctx, table.DefaultTxControl(), query,
			table.NewQueryParameters(
				table.ValueParam("$user_id", types.UTF8Value(userID)),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to remove connection: %w", err)
		}

		return nil
	})
}

// RemoveConnectionByID removes a connection by connection ID
func RemoveConnectionByID(ctx context.Context, db *ydb.Driver, connectionID string) error {
	return db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := `
			DECLARE $connection_id AS Utf8;
			DELETE FROM connections
			WHERE connection_id = $connection_id;
		`

		_, _, err := s.Execute(ctx, table.DefaultTxControl(), query,
			table.NewQueryParameters(
				table.ValueParam("$connection_id", types.UTF8Value(connectionID)),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to remove connection by ID: %w", err)
		}

		return nil
	})
}

// GetAllConnections retrieves all active connections from YDB
func GetAllConnections(ctx context.Context, db *ydb.Driver) ([]Connection, error) {
	var connections []Connection

	err := db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := `
			SELECT connection_id, user_id, connected_at
			FROM connections;
		`

		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query, table.NewQueryParameters())
		if err != nil {
			return fmt.Errorf("failed to query connections: %w", err)
		}
		defer func() {
			_ = res.Close()
		}()

		for res.NextResultSet(ctx) {
			for res.NextRow() {
				var conn Connection
				if err := res.ScanWithDefaults(&conn.ConnectionID, &conn.UserID, &conn.ConnectedAt); err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
				connections = append(connections, conn)
			}
		}

		if err := res.Err(); err != nil {
			return fmt.Errorf("result set error: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return connections, nil
}

// GetUserIDByConnectionID retrieves the user ID for a given connection ID
func GetUserIDByConnectionID(ctx context.Context, db *ydb.Driver, connectionID string) (string, error) {
	var userID string

	err := db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := `
			DECLARE $connection_id AS Utf8;
			SELECT user_id
			FROM connections
			WHERE connection_id = $connection_id;
		`

		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query,
			table.NewQueryParameters(
				table.ValueParam("$connection_id", types.UTF8Value(connectionID)),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to query user ID: %w", err)
		}
		defer func() {
			_ = res.Close()
		}()

		for res.NextResultSet(ctx) {
			if res.NextRow() {
				if err := res.ScanWithDefaults(&userID); err != nil {
					return fmt.Errorf("failed to scan user ID: %w", err)
				}
			}
		}

		if err := res.Err(); err != nil {
			return fmt.Errorf("result set error: %w", err)
		}

		if userID == "" {
			return fmt.Errorf("connection not found: %s", connectionID)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return userID, nil
}
