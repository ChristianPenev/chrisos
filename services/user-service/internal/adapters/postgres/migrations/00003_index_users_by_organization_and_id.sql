-- +goose Up
CREATE INDEX users_organization_id_id_idx ON users (organization_id, id);
DROP INDEX users_organization_id_idx;

-- +goose Down
CREATE INDEX users_organization_id_idx ON users (organization_id);
DROP INDEX users_organization_id_id_idx;
