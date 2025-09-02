-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_log (
    id SERIAL PRIMARY KEY,
    action VARCHAR(255) NOT NULL,
    actor VARCHAR(225) NOT NULL,
    payload JSONB NOT NULL,
    ts TIMESTAMP default NOW()
);

CREATE INDEX ids_audit_log_action ON audit_log(action);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXIST audit_log;
-- +goose StatementEnd
