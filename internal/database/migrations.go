package database

type Migration struct {
	name string
	sql  string
}

// getMigrations returns all database migrations in order
func getMigrations() []Migration {
	return []Migration{
		{
			name: "001_create_projects_table",
			sql: `
				CREATE TABLE IF NOT EXISTS projects (
					id TEXT PRIMARY KEY,
					name TEXT NOT NULL,
					description TEXT,
					color TEXT NOT NULL,
					created_at DATETIME NOT NULL,
					updated_at DATETIME NOT NULL
				);
			`,
		},
		{
			name: "002_create_tasks_table",
			sql: `
				CREATE TABLE IF NOT EXISTS tasks (
					id TEXT PRIMARY KEY,
					project_id TEXT NOT NULL,
					name TEXT NOT NULL,
					desc TEXT,
					status INTEGER NOT NULL,
					priority INTEGER DEFAULT 1,
					created_at DATETIME NOT NULL,
					updated_at DATETIME NOT NULL,
					FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
				);
			`,
		},

		{
			name: "005_create_indexes",
			sql: `
				CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
				CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
				CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
				CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at);
				CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
			`,
		},
	}
}
