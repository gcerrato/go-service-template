package services

import (
	"context"
	"log/slog"

	"github.com/gcerrato/go-service-template/api/models"
	"github.com/gcerrato/go-service-template/database/ent"
	"github.com/gcerrato/go-service-template/internal/repos"
	"github.com/google/uuid"
)

type TodoService struct {
	TodoRepo repos.TodoCRUD
}

func NewTodoService(todoRepo repos.TodoCRUD) *TodoService {
	return &TodoService{TodoRepo: todoRepo}
}

type TodoServiceInterface interface {
	TodoServiceCreator
	TodoServiceGetter
	TodoServiceUpdater
	TodoServiceDeleter
}

type TodoServiceCreator interface {
	CreateTodo(ctx context.Context, newTodo models.TodoCreate) (*ent.Todo, error)
}

type TodoServiceGetter interface {
	GetTodos(ctx context.Context, params models.GetTodosParams) ([]*ent.Todo, error)
	GetTodoById(ctx context.Context, id uuid.UUID) (*ent.Todo, error)
}

type TodoServiceUpdater interface {
	UpdateTodo(ctx context.Context, id uuid.UUID, todo models.TodoUpdate) (*ent.Todo, error)
}

type TodoServiceDeleter interface {
	DeleteTodo(ctx context.Context, id uuid.UUID) error
}

func (t *TodoService) CreateTodo(ctx context.Context, newTodo models.TodoCreate) (*ent.Todo, error) {
	todo, err := t.TodoRepo.CreateTodo(ctx, newTodo)
	if err != nil {
		slog.Error("service error creating todo")
		return nil, err
	}
	return todo, nil
}

func (t *TodoService) GetTodos(ctx context.Context, params models.GetTodosParams) ([]*ent.Todo, error) {
	return t.TodoRepo.GetTodos(ctx, params)
}

func (t *TodoService) GetTodoById(ctx context.Context, id uuid.UUID) (*ent.Todo, error) {
	return t.TodoRepo.GetTodoById(ctx, id)
}

func (t *TodoService) UpdateTodo(ctx context.Context, id uuid.UUID, todo models.TodoUpdate) (*ent.Todo, error) {
	return t.TodoRepo.UpdateTodo(ctx, id, todo)
}

func (t *TodoService) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	return t.TodoRepo.DeleteTodo(ctx, id)
}
