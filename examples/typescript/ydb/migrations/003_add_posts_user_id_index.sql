-- +goose Up
ALTER TABLE posts
    ADD INDEX idx_user_id GLOBAL ON (user_id);

-- +goose Down
ALTER TABLE posts
    DROP INDEX idx_user_id;
