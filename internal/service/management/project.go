package management

import (
	"context"
	"project-management/internal/domain"
	"project-management/internal/domain/project"
	"project-management/pkg/log"
)

func (s *Service) CreateProject(ctx context.Context, req project.Request) (id string, err error) {
	logger := log.LoggerFromContext(ctx)

	data := project.Entity{
		ID:          domain.GenerateID(),
		Title:       req.Title,
		Description: req.Description,
		ManagerID:   req.ManagerID,
		StartedAt:   domain.OnlyDate(req.StartedAt),
		FinishedAt:  domain.OnlyDate(req.FinishedAt),
	}

	id, err = s.projectRepository.Create(ctx, data)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to create project")
		return
	}

	return
}

func (s *Service) GetProject(ctx context.Context, id string) (p project.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	data, err := s.projectRepository.Get(ctx, id)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to get project")
		return
	}

	p = project.ParseFromEntity(data)

	return
}

func (s *Service) UpdateProject(ctx context.Context, id string, req project.UpdateRequest) (err error) {
	logger := log.LoggerFromContext(ctx)

	data := project.Entity{
		Title:       req.Title,
		Description: req.Description,
		ManagerID:   req.ManagerID,
		FinishedAt:  domain.OnlyDate(req.FinishedAt),
	}

	err = s.projectRepository.Update(ctx, id, data)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to update project")
		return
	}

	return
}

func (s *Service) DeleteProject(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx)

	err = s.projectRepository.Delete(ctx, id)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to delete project")
		return
	}

	return
}

func (s *Service) ListProjects(ctx context.Context) (res []project.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	data, err := s.projectRepository.List(ctx)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to list projects")
		return nil, project.ErrNotFound
	}

	res = project.ParseFromEntities(data)

	return
}

func (s *Service) SearchProjects(ctx context.Context, filter, value string) (res []project.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	if value == "" || !project.IsValidFilter(filter) {
		err = project.ErrSearch
		logger.Err(err).Stack().Msg("failed to search tasks")
		return
	}

	data, err := s.projectRepository.Search(ctx, filter, value)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to search projects")
		return
	}

	res = project.ParseFromEntities(data)

	return
}
