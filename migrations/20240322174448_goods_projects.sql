-- +goose Up
-- +goose StatementBegin
CREATE TABLE project (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL
);

INSERT INTO project(name, created_at)
VALUES ('Первая запись', NOW());

CREATE TABLE good (
    id SERIAL PRIMARY KEY,
    project_id SERIAL REFERENCES project (id),
    name VARCHAR NOT NULL,
    description VARCHAR NULL,
    priority INTEGER NOT NULL,
    removed BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE good;
DROP TABLE project;
-- +goose StatementEnd
