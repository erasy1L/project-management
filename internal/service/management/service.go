package management

import (
	"project-management/internal/domain/project"
	"project-management/internal/domain/task"
	"project-management/internal/domain/user"
)

type Service struct {
	userRepostitory   user.Repository
	taskRepository    task.Repository
	projectRepository project.Repository
}

type Configuration func(s *Service) error

func New(cfg ...Configuration) *Service {
	s := &Service{}

	for _, cfg := range cfg {
		cfg(s)
	}

	return s
}

func WithUserRepository(userRepository user.Repository) Configuration {
	return func(s *Service) error {
		s.userRepostitory = userRepository
		return nil
	}
}

func WithTaskRepository(taskRepository task.Repository) Configuration {
	return func(s *Service) error {
		s.taskRepository = taskRepository
		return nil
	}
}

func WithProjectRepository(projectRepository project.Repository) Configuration {
	return func(s *Service) error {
		s.projectRepository = projectRepository
		return nil
	}
}
