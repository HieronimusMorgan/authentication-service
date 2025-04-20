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
    archived_enabled        BOOLEAN   DEFAULT FALSE,
    group_invite_disallowed INT[] DEFAULT NULL,
    group_invite_type       INT       DEFAULT 1,
    group_invite_disallowed INT[] DEFAULT NULL,
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_settings_user_id ON user_settings (user_id);

INSERT INTO user_settings (user_id,
                           archived_enabled,
                           group_invite_type,
                           group_invite_disallowed)
SELECT user_id,
       FALSE,
       1,
       ARRAY['none']
FROM users
WHERE user_id NOT IN (SELECT user_id FROM user_settings);

-- CREATE TABLE family
-- (
--     family_id   SERIAL PRIMARY KEY,
--     family_name VARCHAR(255) NOT NULL,
--     owner_id    INT          NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
--     created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     created_by  VARCHAR(255),
--     updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_by  VARCHAR(255),
--     deleted_at  TIMESTAMP NULL,
--     deleted_by  VARCHAR(255)
-- );
--
-- CREATE TABLE family_invitation_status
-- (
--     status_id   SERIAL PRIMARY KEY,
--     status_name VARCHAR(50) NOT NULL UNIQUE,
--     created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     created_by  VARCHAR(255),
--     updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_by  VARCHAR(255),
--     deleted_at  TIMESTAMP NULL,
--     deleted_by  VARCHAR(255)
-- );
--
-- -- Insert default invitation statuses
-- INSERT INTO family_invitation_status (status_name)
-- VALUES ('Pending'),
--        ('Accepted'),
--        ('Rejected');
--
-- CREATE TABLE family_invitation
-- (
--     invitation_id    SERIAL PRIMARY KEY,
--     family_id        INT NOT NULL REFERENCES family (family_id) ON DELETE CASCADE,
--     sender_user_id   INT NOT NULL REFERENCES users (user_id),
--     receiver_user_id INT NOT NULL REFERENCES users (user_id),
--     status_id        INT NOT NULL REFERENCES family_invitation_status (status_id),
--     invited_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     responded_at     TIMESTAMP DEFAULT NULL,
--     UNIQUE (family_id, receiver_user_id)
-- );
--
-- CREATE INDEX idx_family_invitation_sender_user_id ON family_invitation (sender_user_id);
-- CREATE INDEX idx_family_invitation_receiver_user_id ON family_invitation (receiver_user_id);
-- CREATE INDEX idx_family_invitation_status_id ON family_invitation (status_id);
--
-- CREATE TABLE family_permission
-- (
--     permission_id   SERIAL PRIMARY KEY,
--     permission_name VARCHAR(100) UNIQUE NOT NULL,
--     description     TEXT,
--     created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     created_by      VARCHAR(255),
--     updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_by      VARCHAR(255),
--     deleted_at      TIMESTAMP NULL,
--     deleted_by      VARCHAR(255)
-- );
--
-- -- Insert default family_permission
-- INSERT INTO family_permission (permission_name, description)
-- VALUES ('Admin', 'Full control over family members and assets'),
--        ('Manage', 'Manage family members and permissions'),
--        ('Read-Write', 'Read and Write access to assets'),
--        ('Read', 'Read/View assets');
--
-- CREATE TABLE family_member_permission
-- (
--     family_id     INT REFERENCES family (family_id) ON DELETE CASCADE,
--     user_id       INT REFERENCES users (user_id) ON DELETE CASCADE,
--     permission_id INT REFERENCES family_permission (permission_id) ON DELETE CASCADE,
--     created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     created_by    VARCHAR(255),
--     PRIMARY KEY (family_id, user_id, permission_id)
-- );
--
-- CREATE INDEX idx_fmp_user ON family_member_permission (user_id);
-- CREATE INDEX idx_fmp_family ON family_member_permission (family_id);
-- CREATE INDEX idx_fmp_permission ON family_member_permission (permission_id);
--
-- CREATE TABLE family_member
-- (
--     family_id  INT REFERENCES family (family_id) ON DELETE CASCADE,
--     user_id    INT REFERENCES users (user_id) ON DELETE CASCADE,
--     joined_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     created_by VARCHAR(255),
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_by VARCHAR(255),
--     deleted_at TIMESTAMP NULL,
--     deleted_by VARCHAR(255),
--     PRIMARY KEY (family_id, user_id)
-- );
--
-- CREATE INDEX idx_family_member_family_id ON family_member (family_id);
-- CREATE INDEX idx_family_member_user_id ON family_member (user_id);

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
VALUES ('resource', 'Description of the resource', CURRENT_TIMESTAMP, 'system'),
       ('system', 'Operations for managing the system', CURRENT_TIMESTAMP, 'system'),
       ('auth', 'Authentication-related operations', CURRENT_TIMESTAMP, 'system'),
       ('asset', 'Asset management operations', CURRENT_TIMESTAMP, 'system');

INSERT INTO role_resources (role_id, resource_id, created_at, created_by)
SELECT (SELECT role_id FROM roles WHERE name = 'Super Admin') AS role_id,
       resource_id,
       CURRENT_TIMESTAMP,
       'system'
FROM resources;

INSERT INTO role_resources (role_id, resource_id, created_at, created_by)
VALUES ((SELECT role_id FROM roles WHERE name = 'Admin'),
        (SELECT resource_id FROM resources WHERE name = 'resource'),
        CURRENT_TIMESTAMP, 'system'),
       ((SELECT role_id FROM roles WHERE name = 'Admin'), (SELECT resource_id FROM resources WHERE name = 'auth'),
        CURRENT_TIMESTAMP, 'system'),
       ((SELECT role_id FROM roles WHERE name = 'Admin'), (SELECT resource_id FROM resources WHERE name = 'master'),
        CURRENT_TIMESTAMP, 'system'),
       ((SELECT role_id FROM roles WHERE name = 'Admin'), (SELECT resource_id FROM resources WHERE name = 'asset'),
        CURRENT_TIMESTAMP, 'system');

INSERT INTO role_resources (role_id, resource_id, created_at, created_by)
VALUES ((SELECT role_id FROM roles WHERE name = 'User'), (SELECT resource_id FROM resources WHERE name = 'auth'),
        CURRENT_TIMESTAMP, 'system'),
       ((SELECT role_id FROM roles WHERE name = 'User'), (SELECT resource_id FROM resources WHERE name = 'asset'),
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