package repos

import (
	"context"

	"github.com/gcerrato/go-service-template/api/models"
	"github.com/gcerrato/go-service-template/database/ent"
	"github.com/gcerrato/go-service-template/database/ent/todo"
	"github.com/gcerrato/go-service-template/pkg"
	"github.com/google/uuid"
)

type TodoRepo struct {
	db *ent.Client
}

type TodoCRUD interface {
	TodoCreater
	TodoReader
	TodoUpdater
	TodoDeleter
}

type TodoCreater interface {
	CreateTodo(ctx context.Context, apiModel models.TodoCreate) (*ent.Todo, error)
}

type TodoReader interface {
	GetTodos(ctx context.Context, params models.GetTodosParams) ([]*ent.Todo, error)
	GetTodoById(ctx context.Context, id uuid.UUID) (*ent.Todo, error)
}

type TodoUpdater interface {
	UpdateTodo(ctx context.Context, id uuid.UUID, updateTodo models.TodoUpdate) (*ent.Todo, error)
}

type TodoDeleter interface {
	DeleteTodo(ctx context.Context, id uuid.UUID) error
}

func NewTodoRepo(db *ent.Client) *TodoRepo {
	return &TodoRepo{db: db}
}

func (t *TodoRepo) GetTodos(ctx context.Context, params models.GetTodosParams) ([]*ent.Todo, error) {
	query := t.db.Todo.Query()

	if params.Completed != nil {
		query = query.Where(todo.Completed(*params.Completed))
	}

	if params.Priority != nil {
		query = query.Where(todo.PriorityEQ(todo.Priority(*params.Priority)))
	}

	result, err := query.All(ctx)
	if err != nil {
		return nil, pkg.SendRepoError("repo error querying todos", err)
	}
	return result, nil
}

func (t *TodoRepo) GetTodoById(ctx context.Context, id uuid.UUID) (*ent.Todo, error) {
	result, err := t.db.Todo.Get(ctx, id)
	if err != nil {
		return nil, pkg.SendRepoError("repo error on get todo", err)
	}
	return result, nil
}

func (t *TodoRepo) CreateTodo(ctx context.Context, newTodo models.TodoCreate) (*ent.Todo, error) {
	todoBuilder := t.db.Todo.Create()
	todo, err := todoBuilder.
		SetTitle(newTodo.Title).
		SetCompleted(false).
		SetNillableDescription(newTodo.Description).
		SetNillablePriority((*todo.Priority)(newTodo.Priority)).
		SetNillableDueDate(newTodo.DueDate).
		Save(ctx)

	if err != nil {
		return nil, pkg.SendRepoError("repo error creating todo", err)
	}
	return todo, nil
}

func (t *TodoRepo) UpdateTodo(ctx context.Context, id uuid.UUID, updateTodo models.TodoUpdate) (*ent.Todo, error) {
	todoBuilder := t.db.Todo.UpdateOneID(id)

	if updateTodo.Title != nil {
		todoBuilder.SetTitle(*updateTodo.Title)
	}
	if updateTodo.Description != nil {
		todoBuilder.SetDescription(*updateTodo.Description)
	}
	if updateTodo.Completed != nil {
		todoBuilder.SetCompleted(*updateTodo.Completed)
	}
	if updateTodo.Priority != nil {
		todoBuilder.SetPriority(todo.Priority(*updateTodo.Priority))
	}
	if updateTodo.DueDate != nil {
		todoBuilder.SetDueDate(*updateTodo.DueDate)
	}

	todo, err := todoBuilder.Save(ctx)
	if err != nil {
		return nil, pkg.SendRepoError("repo error updating todo", err)
	}
	return todo, nil
}

func (t *TodoRepo) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	err := t.db.Todo.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return pkg.SendRepoError("repo error deleting todo", err)
	}
	return nil
}
