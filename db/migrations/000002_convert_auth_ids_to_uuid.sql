-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE users ADD COLUMN id_uuid UUID;
UPDATE users
SET id_uuid = CASE
    WHEN id ~ '^[0-9a-fA-F]{32}$' THEN (
        substring(id from 1 for 8) || '-' ||
        substring(id from 9 for 4) || '-' ||
        substring(id from 13 for 4) || '-' ||
        substring(id from 17 for 4) || '-' ||
        substring(id from 21 for 12)
    )::UUID
    ELSE gen_random_uuid()
END
WHERE id_uuid IS NULL;
ALTER TABLE users ALTER COLUMN id_uuid SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_id_uuid_key UNIQUE (id_uuid);

ALTER TABLE auth_password_credentials ADD COLUMN user_id_uuid UUID;
UPDATE auth_password_credentials credential
SET user_id_uuid = users.id_uuid
FROM users
WHERE credential.user_id = users.id;
ALTER TABLE auth_password_credentials ALTER COLUMN user_id_uuid SET NOT NULL;
ALTER TABLE auth_password_credentials ADD CONSTRAINT auth_password_credentials_user_id_uuid_key UNIQUE (user_id_uuid);

ALTER TABLE auth_sessions ADD COLUMN id_uuid UUID;
UPDATE auth_sessions
SET id_uuid = CASE
    WHEN id ~ '^[0-9a-fA-F]{32}$' THEN (
        substring(id from 1 for 8) || '-' ||
        substring(id from 9 for 4) || '-' ||
        substring(id from 13 for 4) || '-' ||
        substring(id from 17 for 4) || '-' ||
        substring(id from 21 for 12)
    )::UUID
    ELSE gen_random_uuid()
END
WHERE id_uuid IS NULL;
ALTER TABLE auth_sessions ALTER COLUMN id_uuid SET NOT NULL;
ALTER TABLE auth_sessions ADD CONSTRAINT auth_sessions_id_uuid_key UNIQUE (id_uuid);

ALTER TABLE auth_sessions ADD COLUMN user_id_uuid UUID;
UPDATE auth_sessions session
SET user_id_uuid = users.id_uuid
FROM users
WHERE session.user_id = users.id;
ALTER TABLE auth_sessions ALTER COLUMN user_id_uuid SET NOT NULL;

DROP INDEX IF EXISTS idx_auth_sessions_user_id;

ALTER TABLE auth_password_credentials DROP CONSTRAINT IF EXISTS auth_password_credentials_user_id_fkey;
ALTER TABLE auth_sessions DROP CONSTRAINT IF EXISTS auth_sessions_user_id_fkey;
ALTER TABLE auth_password_credentials DROP CONSTRAINT IF EXISTS auth_password_credentials_pkey;
ALTER TABLE auth_sessions DROP CONSTRAINT IF EXISTS auth_sessions_pkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey;

ALTER TABLE auth_password_credentials DROP COLUMN user_id;
ALTER TABLE auth_sessions DROP COLUMN user_id;
ALTER TABLE auth_sessions DROP COLUMN id;
ALTER TABLE users DROP COLUMN id;

ALTER TABLE users RENAME COLUMN id_uuid TO id;
ALTER TABLE auth_password_credentials RENAME COLUMN user_id_uuid TO user_id;
ALTER TABLE auth_sessions RENAME COLUMN id_uuid TO id;
ALTER TABLE auth_sessions RENAME COLUMN user_id_uuid TO user_id;

ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE auth_password_credentials ADD CONSTRAINT auth_password_credentials_pkey PRIMARY KEY (user_id);
ALTER TABLE auth_sessions ADD CONSTRAINT auth_sessions_pkey PRIMARY KEY (id);

ALTER TABLE auth_password_credentials
    ADD CONSTRAINT auth_password_credentials_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE auth_sessions
    ADD CONSTRAINT auth_sessions_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id ON auth_sessions(user_id);

-- +goose Down
ALTER TABLE users ADD COLUMN id_text TEXT;
UPDATE users SET id_text = replace(id::TEXT, '-', '');
ALTER TABLE users ALTER COLUMN id_text SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_id_text_key UNIQUE (id_text);

ALTER TABLE auth_password_credentials ADD COLUMN user_id_text TEXT;
UPDATE auth_password_credentials credential
SET user_id_text = replace(users.id::TEXT, '-', '')
FROM users
WHERE credential.user_id = users.id;
ALTER TABLE auth_password_credentials ALTER COLUMN user_id_text SET NOT NULL;
ALTER TABLE auth_password_credentials ADD CONSTRAINT auth_password_credentials_user_id_text_key UNIQUE (user_id_text);

ALTER TABLE auth_sessions ADD COLUMN id_text TEXT;
UPDATE auth_sessions SET id_text = replace(id::TEXT, '-', '');
ALTER TABLE auth_sessions ALTER COLUMN id_text SET NOT NULL;
ALTER TABLE auth_sessions ADD CONSTRAINT auth_sessions_id_text_key UNIQUE (id_text);

ALTER TABLE auth_sessions ADD COLUMN user_id_text TEXT;
UPDATE auth_sessions session
SET user_id_text = replace(users.id::TEXT, '-', '')
FROM users
WHERE session.user_id = users.id;
ALTER TABLE auth_sessions ALTER COLUMN user_id_text SET NOT NULL;

DROP INDEX IF EXISTS idx_auth_sessions_user_id;

ALTER TABLE auth_password_credentials DROP CONSTRAINT IF EXISTS auth_password_credentials_user_id_fkey;
ALTER TABLE auth_sessions DROP CONSTRAINT IF EXISTS auth_sessions_user_id_fkey;
ALTER TABLE auth_password_credentials DROP CONSTRAINT IF EXISTS auth_password_credentials_pkey;
ALTER TABLE auth_sessions DROP CONSTRAINT IF EXISTS auth_sessions_pkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey;

ALTER TABLE auth_password_credentials DROP COLUMN user_id;
ALTER TABLE auth_sessions DROP COLUMN user_id;
ALTER TABLE auth_sessions DROP COLUMN id;
ALTER TABLE users DROP COLUMN id;

ALTER TABLE users RENAME COLUMN id_text TO id;
ALTER TABLE auth_password_credentials RENAME COLUMN user_id_text TO user_id;
ALTER TABLE auth_sessions RENAME COLUMN id_text TO id;
ALTER TABLE auth_sessions RENAME COLUMN user_id_text TO user_id;

ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE auth_password_credentials ADD CONSTRAINT auth_password_credentials_pkey PRIMARY KEY (user_id);
ALTER TABLE auth_sessions ADD CONSTRAINT auth_sessions_pkey PRIMARY KEY (id);

ALTER TABLE auth_password_credentials
    ADD CONSTRAINT auth_password_credentials_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE auth_sessions
    ADD CONSTRAINT auth_sessions_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id ON auth_sessions(user_id);
