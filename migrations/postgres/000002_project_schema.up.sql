CREATE TABLE IF NOT EXISTS projects (
	id VARCHAR(24) PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	description VARCHAR(200) NOT NULL,
	started_at DATE NOT NULL,
	finished_at DATE NOT NULL,
	manager_id VARCHAR(24) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS projects_title_idx ON projects(title);
CREATE INDEX IF NOT EXISTS projects_manager_idx ON projects(manager_id);
