CREATE TABLE IF NOT EXISTS tasks (
	id VARCHAR(24) PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	description VARCHAR(200) NOT NULL,
	priority VARCHAR CHECK (priority IN ('low', 'medium', 'high')),
	status VARCHAR CHECK (status IN ('active', 'in_proccess', 'done')),
	author_id VARCHAR(24) REFERENCES users(id) ON DELETE SET NULL,
	project_id VARCHAR(24) REFERENCES projects(id) ON DELETE CASCADE,
	created_at DATE NOT NULL,
	done_at DATE NOT NULL
);

CREATE INDEX IF NOT EXISTS tasks_title_idx ON tasks(title);
CREATE INDEX IF NOT EXISTS tasks_author_idx ON tasks(author_id);
CREATE INDEX IF NOT EXISTS tasks_project_idx ON tasks(project_id);
CREATE INDEX IF NOT EXISTS tasks_status_idx ON tasks(status);
CREATE INDEX IF NOT EXISTS tasks_priority_idx ON tasks(priority);
