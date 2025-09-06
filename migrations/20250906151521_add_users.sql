-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,
    roles VARCHAR(255)[],
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);

-- username: admin password: rtc
INSERT INTO users (username, password_hash, roles) VALUES ('admin', '$2a$12$NpjWD0uJCuQ61/hLWO9w6.QLKdkT.06Z1EedRWwQoRUksH40OK54e', '{"admin"}');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
