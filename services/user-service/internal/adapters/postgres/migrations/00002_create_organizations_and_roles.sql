-- +goose Up
CREATE TABLE organizations (
    id         uuid        PRIMARY KEY,
    name       text        NOT NULL,
    created_at timestamptz NOT NULL,
    CONSTRAINT organizations_name_not_empty CHECK (name <> '')
);

CREATE TABLE roles (
    id              uuid        PRIMARY KEY,
    organization_id uuid        NOT NULL REFERENCES organizations (id),
    name            text        NOT NULL,
    owner           boolean     NOT NULL DEFAULT false,
    created_at      timestamptz NOT NULL,
    CONSTRAINT roles_organization_id_name_key UNIQUE (organization_id, name),
    CONSTRAINT roles_id_organization_id_key UNIQUE (id, organization_id),
    CONSTRAINT roles_name_not_empty CHECK (name <> '')
);

CREATE UNIQUE INDEX roles_one_owner_per_organization ON roles (organization_id) WHERE owner;

CREATE TABLE role_permissions (
    role_id    uuid NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission text NOT NULL,
    PRIMARY KEY (role_id, permission)
);

ALTER TABLE users
    ADD COLUMN organization_id uuid NOT NULL REFERENCES organizations (id),
    ADD COLUMN role_id         uuid NOT NULL,
    ADD CONSTRAINT users_role_same_organization
        FOREIGN KEY (role_id, organization_id) REFERENCES roles (id, organization_id);

CREATE INDEX users_organization_id_idx ON users (organization_id);
CREATE INDEX users_role_id_organization_id_idx ON users (role_id, organization_id);

-- +goose Down
ALTER TABLE users
    DROP CONSTRAINT users_role_same_organization,
    DROP COLUMN role_id,
    DROP COLUMN organization_id;

DROP TABLE role_permissions;
DROP TABLE roles;
DROP TABLE organizations;
