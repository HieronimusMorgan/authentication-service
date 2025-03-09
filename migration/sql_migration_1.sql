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
    client_id        VARCHAR(255) UNIQUE NOT NULL,
    username         VARCHAR(255) UNIQUE NOT NULL,
    password         TEXT                NOT NULL,
    email            VARCHAR(255) UNIQUE NOT NULL,
    pin_code         TEXT      DEFAULT NULL,
    pin_attempts     INT       DEFAULT 0 CHECK (pin_attempts >= 0),
    pin_last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    first_name       VARCHAR(255),
    last_name        VARCHAR(255),
    full_name        VARCHAR(255),
    phone_number     VARCHAR(50) UNIQUE,
    profile_picture  TEXT,
    role_id          INT                 NOT NULL REFERENCES roles (role_id) ON DELETE CASCADE,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by       VARCHAR(255),
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP NULL,
    deleted_by       VARCHAR(255)
);

CREATE INDEX idx_users_role_id ON users (role_id);
CREATE INDEX idx_users_client_id ON users (client_id);

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

-- RoleResource Table
CREATE TABLE role_resources
(
    role_id     INT NOT NULL REFERENCES roles (role_id) ON DELETE CASCADE,
    resource_id INT NOT NULL REFERENCES resources (resource_id) ON DELETE CASCADE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255),
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP NULL,
    deleted_by  VARCHAR(255),
    PRIMARY KEY (role_id, resource_id)
);

CREATE INDEX idx_role_resources_role_id ON role_resources (role_id);
CREATE INDEX idx_role_resources_resource_id ON role_resources (resource_id);

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
VALUES ('New Resource', 'Description of the new resource', CURRENT_TIMESTAMP, 'system'),
       ('System Management', 'Operations for managing the entire system', CURRENT_TIMESTAMP, 'system'),
       ('Auth', 'Authentication-related operations', CURRENT_TIMESTAMP, 'system'),
       ('Master', 'Master data management operations', CURRENT_TIMESTAMP, 'system');

INSERT INTO role_resources (role_id, resource_id, created_at, created_by)
SELECT (SELECT role_id FROM roles WHERE name = 'Super Admin') AS role_id,
       resource_id,
       CURRENT_TIMESTAMP,
       'system'
FROM resources;

INSERT INTO role_resources (role_id, resource_id, created_at, created_by)
VALUES ((SELECT role_id FROM roles WHERE name = 'Admin'),
        (SELECT resource_id FROM resources WHERE name = 'New Resource'),
        CURRENT_TIMESTAMP, 'system'),
       ((SELECT role_id FROM roles WHERE name = 'Admin'), (SELECT resource_id FROM resources WHERE name = 'Auth'),
        CURRENT_TIMESTAMP, 'system'),
       ((SELECT role_id FROM roles WHERE name = 'Admin'), (SELECT resource_id FROM resources WHERE name = 'Master'),
        CURRENT_TIMESTAMP, 'system');

INSERT INTO role_resources (role_id, resource_id, created_at, created_by)
VALUES ((SELECT role_id FROM roles WHERE name = 'User'), (SELECT resource_id FROM resources WHERE name = 'Auth'),
        CURRENT_TIMESTAMP, 'system');

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
CREATE TRIGGER set_updated_at_role_resources
    BEFORE UPDATE
    ON role_resources
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
