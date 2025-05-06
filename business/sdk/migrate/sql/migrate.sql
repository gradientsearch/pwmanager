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