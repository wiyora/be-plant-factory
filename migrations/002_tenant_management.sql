CREATE TYPE tenant_status AS ENUM('active', 'inactive', 'suspended');

CREATE TABLE tenants (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	name VARCHAR(32) NOT NULL,
	slug VARCHAR(64) UNIQUE NOT NULL,
	status tenant_status DEFAULT 'active' NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ
);

---- create above / drop below ----

DROP TABLE IF EXISTS tenants;
DROP TYPE IF EXISTS tenant_status;
