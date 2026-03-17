-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS base_role TEXT;

UPDATE users
SET base_role = 'user'
WHERE base_role IS NULL;

ALTER TABLE users
    ALTER COLUMN base_role SET NOT NULL;

ALTER TABLE auth_sessions
    ADD COLUMN IF NOT EXISTS replaced_by_session_id UUID NULL REFERENCES auth_sessions(id) ON DELETE SET NULL;

ALTER TABLE auth_sessions
    ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ NULL;

CREATE TABLE IF NOT EXISTS auth_email_verification_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_password_reset_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS access_permissions (
    id UUID PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS access_sub_roles (
    id UUID PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS access_sub_role_permissions (
    sub_role_id UUID NOT NULL REFERENCES access_sub_roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES access_permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (sub_role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS access_user_sub_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sub_role_id UUID NOT NULL REFERENCES access_sub_roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, sub_role_id)
);

CREATE INDEX IF NOT EXISTS idx_auth_email_verification_tokens_user_id ON auth_email_verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_password_reset_tokens_user_id ON auth_password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_access_user_sub_roles_user_id ON access_user_sub_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_access_sub_role_permissions_sub_role_id ON access_sub_role_permissions(sub_role_id);

-- +goose Down
DROP INDEX IF EXISTS idx_access_sub_role_permissions_sub_role_id;
DROP INDEX IF EXISTS idx_access_user_sub_roles_user_id;
DROP INDEX IF EXISTS idx_auth_password_reset_tokens_user_id;
DROP INDEX IF EXISTS idx_auth_email_verification_tokens_user_id;

DROP TABLE IF EXISTS access_user_sub_roles;
DROP TABLE IF EXISTS access_sub_role_permissions;
DROP TABLE IF EXISTS access_sub_roles;
DROP TABLE IF EXISTS access_permissions;
DROP TABLE IF EXISTS auth_password_reset_tokens;
DROP TABLE IF EXISTS auth_email_verification_tokens;

ALTER TABLE auth_sessions
    DROP COLUMN IF EXISTS last_used_at;

ALTER TABLE auth_sessions
    DROP COLUMN IF EXISTS replaced_by_session_id;

ALTER TABLE users
    DROP COLUMN IF EXISTS base_role;
