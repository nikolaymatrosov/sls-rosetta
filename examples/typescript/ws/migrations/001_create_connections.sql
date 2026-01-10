-- +goose Up
CREATE TABLE connections (
    connection_id Utf8 NOT NULL,
    user_id Utf8,
    connected_at Timestamp,
    PRIMARY KEY (connection_id),
    INDEX idx_user_id GLOBAL ON (user_id)
)
WITH (
    TTL = Interval("P1D") ON connected_at
);

-- +goose Down
DROP TABLE connections;
