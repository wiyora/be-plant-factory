CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TYPE user_current_step AS ENUM ('initial', 'completed');
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'banned');

CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	email VARCHAR(255) UNIQUE NOT NULL,
	name VARCHAR(64) NOT NULL,
    avatar VARCHAR(32) NOT NULL DEFAULT '',
    current_step user_current_step NOT NULL DEFAULT 'initial',
    status user_status NOT NULL DEFAULT 'active',
	is_super_admin BOOLEAN DEFAULT FALSE NOT NULL,
    last_logged_in_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ
);

CREATE INDEX idx_users_name_search ON users USING gin (name gin_trgm_ops);

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
DROP TYPE IF EXISTS user_current_step;
DROP TYPE IF EXISTS user_status;
