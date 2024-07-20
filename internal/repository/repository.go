package repository

import (
	"project-management/config"
	"project-management/internal/domain/project"
	"project-management/internal/domain/task"
	"project-management/internal/domain/user"
	"project-management/internal/repository/postgres"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres postgres.DB

	User    user.Repository
	Task    task.Repository
	Project project.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	s = &Repository{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithPostgresStore(cfg config.DB) Configuration {
	return func(s *Repository) (err error) {
		// Create the postgres store, if we needed parameters, such as connection strings they could be inputted here
		s.postgres, err = postgres.New(cfg)
		if err != nil {
			return
		}

		err = s.postgres.Migrate()
		if err != nil {
			return
		}

		s.User = postgres.NewUserRepository(s.postgres.Client)
		s.Task = postgres.NewTaskRepository(s.postgres.Client)
		s.Project = postgres.NewProjectRepository(s.postgres.Client)

		return
	}
}
