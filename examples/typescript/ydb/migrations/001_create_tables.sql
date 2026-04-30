-- +goose Up
CREATE TABLE users (
    id UUID NOT NULL,
    name Utf8,
    email Utf8,
    PRIMARY KEY (id)
);

CREATE TABLE posts (
    id UUID NOT NULL,
    user_id UUID,
    title Utf8,
    content Utf8,
    PRIMARY KEY (id)
);

CREATE TABLE comments (
    id UUID NOT NULL,
    post_id UUID,
    user_id UUID,
    content Utf8,
    PRIMARY KEY (id)
);

CREATE TABLE likes (
    id UUID NOT NULL,
    post_id UUID,
    user_id UUID,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE likes;
DROP TABLE comments;
DROP TABLE posts;
DROP TABLE users;
