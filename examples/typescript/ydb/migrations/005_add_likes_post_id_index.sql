-- +goose Up
ALTER TABLE likes
    ADD INDEX idx_post_id GLOBAL ON (post_id);

-- +goose Down
ALTER TABLE likes
    DROP INDEX idx_post_id;
