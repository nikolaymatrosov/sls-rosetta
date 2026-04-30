-- +goose Up
CREATE TABLE blobs (
    key        Utf8   NOT NULL,
    partNumber Uint8  NOT NULL,
    data       String,
    PRIMARY KEY (key, partNumber)
);

-- +goose Down
DROP TABLE blobs;
