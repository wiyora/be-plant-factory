CREATE TABLE roles (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	name VARCHAR(32) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ
);

CREATE INDEX idx_roles_name_search ON roles USING gin (name gin_trgm_ops);

CREATE TABLE permissions (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	code VARCHAR(32) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE role_permissions (
	role_id UUID REFERENCES roles (id) ON DELETE CASCADE,
	permission_id UUID REFERENCES permissions (id) ON DELETE CASCADE,
	PRIMARY KEY (role_id, permission_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
