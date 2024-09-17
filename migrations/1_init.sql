-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS project (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS task (
    project_id INT,
    task_id INT PRIMARY KEY

);

CREATE TABLE IF NOT EXISTS users (
    project_id INT,
    user_id INT

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS task;
DROP TABLE IF EXISTS project;
-- +goose StatementEnd
