-- +goose Up
ALTER TABLE comments
    ADD INDEX idx_post_id GLOBAL ON (post_id);

-- +goose Down
ALTER TABLE comments
    DROP INDEX idx_post_id;
