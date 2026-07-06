CREATE TABLE roles (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	name VARCHAR(32) UNIQUE NOT NULL,
    total_permission INT NOT NULL DEFAULT 0,
    permissions JSONB DEFAULT '[]'::jsonb NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ
);

CREATE INDEX idx_roles_name_search ON roles USING gin (name gin_trgm_ops);

---- create above / drop below ----

DROP TABLE IF EXISTS roles;
