-- +goose Up
-- +goose StatementBegin
CREATE TABLE projects
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(300),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_projects_name ON projects(name);

CREATE TABLE environments(
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,

    CONSTRAINT fk_environments_projects
        FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,
    UNIQUE (project_id, name)
);

CREATE TABLE releases (
    id SERIAL PRIMARY KEY,
    environment_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_releases_environments
      FOREIGN KEY (environment_id)
      REFERENCES environments(id)
      ON DELETE CASCADE,
    UNIQUE (environment_id, name)
);

CREATE TABLE configs (
    id SERIAL PRIMARY KEY,
    release_id INTEGER NOT NULL,
    key VARCHAR(255) NOT NULL,
    value_type VARCHAR(255) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,

    CONSTRAINT fk_configs_releases
        FOREIGN KEY (release_id)
        REFERENCES releases(id)
        ON DELETE CASCADE,
    UNIQUE (release_id, key)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXIST projects;
DROP TABLE IF EXIST environments;
DROP TABLE IF EXIST releases;
DROP TABLE IF EXIST configs;
-- +goose StatementEnd
