package task

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidName = errors.New("task name is required")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]Task, error) { 
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, name string) (Task, error) { 
	name = strings.TrimSpace(name)
	if name == "" {
		return Task{}, ErrInvalidName
	}
	return s.repo.Create(ctx, name)
}