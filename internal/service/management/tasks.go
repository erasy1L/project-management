package management

import (
	"context"
	"project-management/internal/domain"
	"project-management/internal/domain/task"
	"project-management/pkg/log"
)

func (s *Service) CreateTask(ctx context.Context, req task.Request) (id string, err error) {
	logger := log.LoggerFromContext(ctx)

	data := task.Entity{
		ID:          domain.GenerateID(),
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		CreatedAt:   domain.OnlyDate(req.CreatedAt),
		DoneAt:      domain.OnlyDate(req.DoneAt),
		AuthorID:    req.AuthorID,
		ProjectID:   req.ProjectID,
	}

	id, err = s.taskRepository.Create(ctx, data)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to create task")
		return
	}

	return
}

func (s *Service) GetTask(ctx context.Context, id string) (res task.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	data, err := s.taskRepository.Get(ctx, id)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to get task")
		return
	}

	res = task.ParseFromEntity(data)

	return
}

func (s *Service) UpdateTask(ctx context.Context, id string, req task.UpdateRequest) (err error) {
	logger := log.LoggerFromContext(ctx)

	data := task.Entity{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		DoneAt:      domain.OnlyDate(req.DoneAt),
		AuthorID:    req.AuthorID,
		ProjectID:   req.AuthorID,
	}

	err = s.taskRepository.Update(ctx, id, data)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to update task")
		return
	}

	return
}

func (s *Service) DeleteTask(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx)

	err = s.taskRepository.Delete(ctx, id)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to delete task")
		return
	}

	return
}

func (s *Service) ListTasks(ctx context.Context) (res []task.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	data, err := s.taskRepository.List(ctx)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to get tasks")
		return
	}

	res = task.ParseFromEntities(data)

	return
}

func (s *Service) SearchTasks(ctx context.Context, filter, value string) (res []task.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	if value == "" || !task.IsValidFilter(filter) {
		err = task.ErrSearch
		logger.Err(err).Stack().Msg("failed to search tasks")
		return
	}

	data, err := s.taskRepository.Search(ctx, filter, value)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to search tasks")
		return
	}

	res = task.ParseFromEntities(data)

	return
}
