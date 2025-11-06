package service

import (
	"context"

	"github.com/msayib/todo-fiber-dig/internal/model"
	"github.com/msayib/todo-fiber-dig/internal/repository"
)

type TodoService interface {
	Create(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error)
	GetByID(ctx context.Context, id uint) (model.Todo, error)
	GetAll(ctx context.Context) ([]model.Todo, error)
	Update(ctx context.Context, id uint, req model.UpdateTodoRequest) (model.Todo, error)
	Delete(ctx context.Context, id uint) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo}
}

func (s *todoService) Create(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error) {
	todo := model.Todo{
		Title:       req.Title,
		Description: req.Description,
	}
	return s.repo.Save(ctx, todo)
}

func (s *todoService) GetByID(ctx context.Context, id uint) (model.Todo, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *todoService) GetAll(ctx context.Context) ([]model.Todo, error) {
	return s.repo.FindAll(ctx)
}

func (s *todoService) Update(ctx context.Context, id uint, req model.UpdateTodoRequest) (model.Todo, error) {
	existingTodo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Todo{}, err
	}

	// Timpa dengan data baru jika ada
	if req.Title != "" {
		existingTodo.Title = req.Title
	}
	if req.Description != "" {
		existingTodo.Description = req.Description
	}
	if req.IsDone != nil {
		existingTodo.IsDone = *req.IsDone
	}

	return s.repo.Update(ctx, existingTodo)
}

func (s *todoService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
