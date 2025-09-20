-- This script should be run once to populate the initial roles, permissions,
-- and the mapping between them. It uses the final, granular permission set.

-- Reset sequences to start IDs from 1 (optional, but good for clean setup)
ALTER SEQUENCE roles_id_seq RESTART WITH 1;
ALTER SEQUENCE permissions_id_seq RESTART WITH 1;

-- Create the roles
INSERT INTO roles (name) VALUES ('user'), ('admin') ON CONFLICT (name) DO NOTHING;

-- Create the granular permissions
INSERT INTO permissions (name) VALUES
    ('files:upload'),
    ('files:download'),
    ('files:delete'),
    ('files:read:shared'),
    ('shares:create:public'),
    ('shares:create:user'),
    ('shares:revoke:public'),
    ('shares:revoke:user'),
    ('stats:read:self'),
    ('admin:manage_roles'),
    ('admin:view_stats')
ON CONFLICT (name) DO NOTHING;

-- Map permissions to roles
-- NOTE: The queries below assume standard IDs. Adjust if you have existing data.

-- 'user' role gets standard file and sharing permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES
    (1, 1), -- user can files:upload
    (1, 2), -- user can files:download
    (1, 3), -- user can files:delete
    (1, 4), -- user can files:read:shared
    (1, 5), -- user can shares:create:public
    (1, 6), -- user can shares:create:user
    (1, 7), -- user can shares:revoke:public
    (1, 8), -- user can shares:revoke:user
    (1, 9)  -- user can stats:read:self
ON CONFLICT DO NOTHING;

-- 'admin' role gets ALL permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES
    (2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 6), (2, 7), (2, 8), (2, 9), (2, 10), (2, 11)
ON CONFLICT DO NOTHING;