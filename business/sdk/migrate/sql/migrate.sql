-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
	user_id       UUID        NOT NULL,
	name          TEXT        NOT NULL,
	email         TEXT UNIQUE NOT NULL,
	roles         TEXT[]      NOT NULL,
	password_hash TEXT        NOT NULL,
    department    TEXT        NULL,
    enabled       BOOLEAN     NOT NULL,
	date_created  TIMESTAMP   NOT NULL,
	date_updated  TIMESTAMP   NOT NULL,

	PRIMARY KEY (user_id)
);

-- Version: 1.02
-- Description: Create table keys
CREATE TABLE keys (
	key_id   UUID           NOT NULL,
    user_id      UUID           NOT NULL,
	name         TEXT           NOT NULL,
    cost         NUMERIC(10, 2) NOT NULL,
	quantity     INT            NOT NULL,
	date_created TIMESTAMP      NOT NULL,
	date_updated TIMESTAMP      NOT NULL,

	PRIMARY KEY (key_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.03
-- Description: Add keys view.
CREATE OR REPLACE VIEW view_keys AS
SELECT
    p.key_id,
    p.user_id,
	p.name,
    p.cost,
	p.quantity,
    p.date_created,
    p.date_updated,
    u.name AS user_name
FROM
    keys AS p
JOIN
    users AS u ON u.user_id = p.user_id

-- Version: 1.04
-- Description: Create table bundles
CREATE TABLE bundles (
    bundle_id       UUID       NOT NULL,
    type          TEXT       NOT NULL,
    user_id       UUID       NOT NULL,
    address_1     TEXT       NOT NULL,
    address_2     TEXT       NULL,
    zip_code      TEXT       NOT NULL,
    city          TEXT       NOT NULL,
    state         TEXT       NOT NULL,
    country       TEXT       NOT NULL,
    date_created  TIMESTAMP  NOT NULL,
    date_updated  TIMESTAMP  NOT NULL,

    PRIMARY KEY (bundle_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
