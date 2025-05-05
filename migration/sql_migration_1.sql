-- Roles Table
CREATE TABLE roles
(
    role_id     SERIAL PRIMARY KEY,
    name        VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255),
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP NULL,
    deleted_by  VARCHAR(255)
);

-- Users Table
CREATE TABLE users
(
    user_id          SERIAL PRIMARY KEY,
    client_id VARCHAR(255) UNIQUE NOT NULL,
    username  VARCHAR(255) UNIQUE NOT NULL,
    password  TEXT                NOT NULL,
    email     VARCHAR(255) UNIQUE NOT NULL,
    pin_code         TEXT      DEFAULT NULL,
    pin_attempts     INT       DEFAULT 0 CHECK (pin_attempts >= 0),
    pin_last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    first_name       VARCHAR(255),
    last_name        VARCHAR(255),
    full_name        VARCHAR(255),
    phone_number     VARCHAR(50) UNIQUE,
    profile_picture  TEXT,
    role_id   INT                 NOT NULL,
    device_id VARCHAR(100),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by       VARCHAR(255),
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP NULL,
    deleted_by       VARCHAR(255)
);

CREATE INDEX idx_users_client_id ON users (client_id);
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);

CREATE TABLE user_settings
(
    setting_id              SERIAL PRIMARY KEY,
    user_id                 INT NOT NULL UNIQUE,
    group_invite_type       INT       DEFAULT 1,
    group_invite_disallowed INT[] DEFAULT NULL,
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_settings_user_id ON user_settings (user_id);

INSERT INTO user_settings (user_id,
                           group_invite_type)
SELECT user_id,
       1
FROM users
WHERE user_id NOT IN (SELECT user_id FROM user_settings);

-- Resources Table
CREATE TABLE resources
(
    resource_id SERIAL PRIMARY KEY,
    name        VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255),
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP NULL,
    deleted_by  VARCHAR(255)
);

-- UserRole Table
CREATE TABLE user_roles
(
    user_id    INT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    role_id    INT NOT NULL REFERENCES roles (role_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(255),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);

-- UserResources Table
CREATE TABLE user_resources
(
    user_id     INT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    resource_id INT NOT NULL REFERENCES resources (resource_id) ON DELETE CASCADE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255),
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP NULL,
    deleted_by  VARCHAR(255),
    PRIMARY KEY (user_id, resource_id)
);

CREATE INDEX idx_user_resources_user_id ON user_resources (user_id);
CREATE INDEX idx_user_resources_resource_id ON user_resources (resource_id);

CREATE TABLE internal_tokens
(
    id          SERIAL PRIMARY KEY,
    resource_id INT     NOT NULL REFERENCES resources (resource_id) ON DELETE CASCADE,
    token       TEXT    NOT NULL,
    expired     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255),
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP NULL,
    deleted_by  VARCHAR(255)
);

CREATE INDEX idx_internal_tokens_resource_id ON internal_tokens (resource_id);
CREATE INDEX idx_internal_tokens_deleted_at ON internal_tokens (deleted_at);

CREATE TABLE user_sessions
(
    user_session_id SERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    session_token   TEXT UNIQUE NOT NULL,
    refresh_token   VARCHAR(255) UNIQUE,
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    login_time      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at      TIMESTAMP   NOT NULL,
    logout_time     TIMESTAMP NULL,
    is_active       BOOLEAN   DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(255)
);
CREATE INDEX idx_user_sessions_user_session_id ON user_sessions (user_session_id);
CREATE INDEX idx_user_sessions_user_id ON user_sessions (user_id);
CREATE INDEX idx_user_sessions_session_token ON user_sessions (session_token);

CREATE TABLE user_keys
(
    user_id               INT PRIMARY KEY REFERENCES users (user_id) ON DELETE CASCADE,
    public_key            TEXT NOT NULL,
    encrypted_private_key TEXT NOT NULL,
    encryption_algorithm  VARCHAR(50) DEFAULT 'RSA-2048',
    salt                  TEXT NOT NULL,
    created_at            TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by            VARCHAR(255),
    updated_at            TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_by            VARCHAR(255),
    deleted_at            TIMESTAMP NULL,
    deleted_by            VARCHAR(255)
);
CREATE INDEX idx_user_keys_user_id ON user_keys (user_id);

CREATE TABLE cron_jobs
(
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL UNIQUE,
    schedule VARCHAR(255) NOT NULL,
    is_active        BOOLEAN   DEFAULT TRUE,
    description      TEXT,
    last_executed_at TIMESTAMP,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by       VARCHAR(255),
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP,
    deleted_by       VARCHAR(255)
);

INSERT INTO cron_jobs (name, schedule, is_active, description, created_by)
VALUES ('user_session_cleanup', '0 5 * * *', true, 'Check User Session Expired', 'system'),
       ('reset_pin_attempts', '0 0 * * *', true, 'Description of cron job 2', 'system');

INSERT INTO roles (name, description, created_at, created_by)
VALUES ('Super Admin', 'Super Administrator with highest privileges', CURRENT_TIMESTAMP, 'system'),
       ('Admin', 'Administrator with full access', CURRENT_TIMESTAMP, 'system'),
       ('User', 'Regular user with limited access', CURRENT_TIMESTAMP, 'system');

INSERT INTO resources (name, description, created_at, created_by)
VALUES ('resource', 'Description of the resource', CURRENT_TIMESTAMP, 'system'),
       ('system', 'Operations for managing the system', CURRENT_TIMESTAMP, 'system'),
       ('auth', 'Authentication-related operations', CURRENT_TIMESTAMP, 'system'),
       ('asset', 'Asset management operations', CURRENT_TIMESTAMP, 'system'),
       ('asset-group', 'Asset group management operations', CURRENT_TIMESTAMP, 'system');

INSERT INTO users (client_id, username, password, first_name, last_name, full_name, phone_number, email,
                   profile_picture,
                   role_id, created_by, updated_by)
VALUES ('super-admin-client-id',
        'super_admin',
        '$2b$12$uZN22OA4kvuE7kWS2ArXD./vHOs/x6dGxEH8ZENMVSS7/XRziCcC.',
        'Super',
        'Admin',
        'Super Admin',
        '1234567890',
        'admin@gmail.com',
        'https://example.com/admin.png',
        (SELECT role_id FROM roles WHERE name = 'Super Admin'),
        'system',
        'system');

INSERT INTO user_roles (user_id, role_id, created_at, created_by)
VALUES ((SELECT user_id FROM users WHERE username = 'super_admin'), -- Super Admin user ID
        (SELECT role_id FROM roles WHERE name = 'Super Admin'), -- Super Admin role ID
        CURRENT_TIMESTAMP, -- Current timestamp
        'system' -- Created by system
       );

INSERT INTO user_resources (user_id, resource_id, created_at, created_by)
SELECT u.user_id,
       r.resource_id,
       CURRENT_TIMESTAMP,
       'system'
FROM users u
         JOIN
     resources r ON r.name IN ('resource', 'system', 'auth', 'asset', 'asset-group')
WHERE u.username = 'super_admin';

-- Function to update updated_at column
CREATE
OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at
= CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Triggers to update updated_at column
CREATE TRIGGER set_updated_at_roles
    BEFORE UPDATE
    ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE
    ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_resources
    BEFORE UPDATE
    ON resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_user_resources
    BEFORE UPDATE
    ON user_resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_user_roles
    BEFORE UPDATE
    ON user_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_internal_tokens
    BEFORE UPDATE
    ON internal_tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_user_sessions
    BEFORE UPDATE
    ON user_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_cron_jobs
    BEFORE UPDATE
    ON cron_jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE
OR REPLACE FUNCTION sync_user_settings()
RETURNS TRIGGER AS $$
BEGIN
INSERT INTO user_settings (user_id)
VALUES (NEW.user_id);
RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER trg_sync_user_settings
    AFTER INSERT
    ON users
    FOR EACH ROW
    EXECUTE FUNCTION sync_user_settings();

-- Allow access to the public schema
GRANT
USAGE
ON
SCHEMA
public TO replicator;
-- Allow SELECT on the users table
GRANT SELECT ON TABLE public.users TO replicator;
-- Allow SELECT on the user_settings table
GRANT SELECT ON TABLE public.user_settings TO replicator;