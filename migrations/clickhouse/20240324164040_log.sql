-- +goose Up
-- +goose StatementBegin
CREATE DATABASE IF NOT EXISTS db_log ENGINE = KeeperMap('/goose_version');

CREATE TABLE db_log.LogInfo (
    id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR,
    priority INTEGER NOT NULL,
    removed BOOL NOT NULL,
    EventTime TIMESTAMP NOT NULL
)
ENGINE = KeeperMap('/goose_version')
PRIMARY KEY (id)
ORDER BY(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
