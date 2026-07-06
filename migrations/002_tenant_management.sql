CREATE TYPE tenant_status AS ENUM('active', 'inactive', 'suspended');

CREATE TABLE tenants (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	name VARCHAR(32) UNIQUE NOT NULL,
    logo VARCHAR(16) NOT NULL,
	status tenant_status DEFAULT 'active' NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ
);

CREATE INDEX idx_tenants_name_search ON tenants USING gin (name gin_trgm_ops);

---- create above / drop below ----

DROP TABLE IF EXISTS tenants;
DROP TYPE IF EXISTS tenant_status;
