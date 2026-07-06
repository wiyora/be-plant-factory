CREATE TABLE user_tenants (
	user_id UUID REFERENCES users (id) ON DELETE CASCADE,
	tenant_id UUID REFERENCES tenants (id) ON DELETE CASCADE,
	role_id UUID REFERENCES roles (id),
	PRIMARY KEY (user_id, tenant_id),
    UNIQUE (user_id, tenant_id),
    UNIQUE (user_id, tenant_id, role_id)
);

CREATE INDEX idx_user_tenants_tenant_id ON user_tenants (tenant_id);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_user_tenants_tenant_id;
DROP TABLE IF EXISTS user_tenants;
