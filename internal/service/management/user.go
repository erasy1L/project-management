package management

import (
	"context"
	"project-management/internal/domain"
	"project-management/internal/domain/user"
	"project-management/pkg/log"
)

func (s *Service) ListUsers(ctx context.Context) (res []user.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	data, err := s.userRepostitory.List(ctx)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to get users")
		return
	}

	res = user.ParseFromEntities(data)

	return
}

func (s *Service) CreateUser(ctx context.Context, req user.Request) (id string, err error) {
	logger := log.LoggerFromContext(ctx)

	data := user.Entity{
		ID:               domain.GenerateID(),
		Name:             req.Name,
		Email:            req.Email,
		RegistrationDate: domain.OnlyDate(req.RegistrationDate),
		Role:             req.Role,
	}

	id, err = s.userRepostitory.Create(ctx, data)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to create user")
		return
	}

	return
}

func (s *Service) GetUser(ctx context.Context, id string) (res user.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	data, err := s.userRepostitory.Get(ctx, id)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to get user")
		return
	}

	res = user.ParseFromEntity(data)

	return
}

func (s *Service) UpdateUser(ctx context.Context, id string, req user.UpdateRequest) (err error) {
	logger := log.LoggerFromContext(ctx)

	data := user.Entity{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	err = s.userRepostitory.Update(ctx, id, data)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to update user")
		return
	}

	return
}

func (s *Service) DeleteUser(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx)

	err = s.userRepostitory.Delete(ctx, id)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to delete user")
		return
	}

	return
}

func (s *Service) SearchUsers(ctx context.Context, filter, value string) (res []user.Response, err error) {
	logger := log.LoggerFromContext(ctx)

	if value == "" || !user.IsValidFilter(filter) {
		err = user.ErrSearch
		logger.Err(err).Stack().Msg("failed to search tasks")
		return
	}

	data, err := s.userRepostitory.Search(ctx, filter, value)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to search users")
		return
	}

	res = user.ParseFromEntities(data)

	return
}
