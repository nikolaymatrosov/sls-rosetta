-- Create users table for the YDB example
-- Run this script in your YDB database after deployment

CREATE TABLE users (
    id Int32 NOT NULL,
    name Utf8,
    PRIMARY KEY (id)
);
