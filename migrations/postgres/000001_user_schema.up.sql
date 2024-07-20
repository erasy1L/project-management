CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(24) PRIMARY KEY,
	name VARCHAR NOT NULL,
	email VARCHAR(32) UNIQUE NOT NULL,
	registration_date DATE NOT NULL,
	role VARCHAR CHECK (role IN ('admin', 'manager', 'developer'))
);

CREATE INDEX IF NOT EXISTS users_name_idx ON users(email);