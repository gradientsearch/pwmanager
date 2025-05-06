-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    roles TEXT [] NOT NULL,
    password_hash TEXT NOT NULL,
    department TEXT NULL,
    enabled BOOLEAN NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE bundles (
    bundle_id UUID NOT NULL,
    user_id UUID NOT NULL,
    type TEXT NOT NULL,
    metadata TEXT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    PRIMARY KEY (bundle_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE keys (
    key_id UUID NOT NULL,
    user_id UUID NOT NULL,
    bundle_id UUID NOT NULL,
    data TEXT NOT NULL,
    roles TEXT [] NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, bundle_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (bundle_id) REFERENCES bundles(bundle_id) ON DELETE CASCADE
);

CREATE TABLE entries (
    entry_id UUID NOT NULL,
    user_id UUID NOT NULL,
    bundle_id UUID NOT NULL,
    data TEXT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    PRIMARY KEY (entry_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (bundle_id) REFERENCES bundles(bundle_id) ON DELETE CASCADE
);

CREATE
OR REPLACE VIEW view_keys AS
SELECT
    p.key_id,
    p.user_id,
    p.date_created,
    p.date_updated,
    u.name AS user_name
FROM
    keys AS p
    JOIN users AS u ON u.user_id = p.user_id


-- -- SQL view (static, no user_id in WHERE clause)
-- CREATE VIEW user_bundles_with_keys AS
-- SELECT
--     u.user_id,
--     u.name,
--     u.email,
--     b.bundle_id,
--     b.type AS bundle_type,
--     b.metadata AS bundle_metadata,
--     b.date_created AS bundle_date_created,
--     b.date_updated AS bundle_date_updated,
--     k.key_id,
--     k.data AS key_data,
--     k.roles AS key_roles,
--     k.date_created AS key_date_created,
--     k.date_updated AS key_date_updated,
--     (
--         SELECT json_agg(json_build_object('user_id', ku.user_id, 'name', ku.name, 'email', ku.email))
--         FROM keys k2
--         JOIN users ku ON k2.user_id = ku.user_id
--         WHERE k2.bundle_id = b.bundle_id
--     ) AS users_with_access
-- FROM
--     users u
-- JOIN
--     bundles b ON b.user_id = u.user_id
-- LEFT JOIN
--     keys k ON k.bundle_id = b.bundle_id AND k.user_id = u.user_id;
-- WHERE b.user_id = ''