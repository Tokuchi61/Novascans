CREATE TABLE access_permissions (
    id UUID PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE access_sub_roles (
    id UUID PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE access_sub_role_permissions (
    sub_role_id UUID NOT NULL REFERENCES access_sub_roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES access_permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (sub_role_id, permission_id)
);

CREATE TABLE access_user_sub_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sub_role_id UUID NOT NULL REFERENCES access_sub_roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, sub_role_id)
);

CREATE INDEX idx_access_user_sub_roles_user_id ON access_user_sub_roles(user_id);
CREATE INDEX idx_access_sub_role_permissions_sub_role_id ON access_sub_role_permissions(sub_role_id);
