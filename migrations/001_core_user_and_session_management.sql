CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	email VARCHAR(255) UNIQUE NOT NULL,
	google_id VARCHAR(32) UNIQUE NOT NULL,
	name VARCHAR(64),
	is_super_admin BOOLEAN DEFAULT FALSE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ
);

CREATE TABLE sessions (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	refresh_token_hash VARCHAR(64) NOT NULL,
	device_name VARCHAR(128) NOT NULL DEFAULT '',
	ip_address INET NOT NULL,
	expired_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ
);

CREATE INDEX idx_sessions_refresh_hash ON sessions (refresh_token_hash);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_sessions_refresh_hash;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
