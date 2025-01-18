package handlers

import (
	"net/http"

	"github.com/gcerrato/go-service-template/api/models"
	"github.com/gcerrato/go-service-template/internal/services"
	"github.com/gcerrato/go-service-template/pkg"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TodosAPI interface {
	PostTodos(ctx echo.Context) error
	GetTodos(ctx echo.Context, params models.GetTodosParams) error
	GetTodosId(ctx echo.Context, id string) error
	PutTodosId(ctx echo.Context, id string) error
	DeleteTodosId(ctx echo.Context, id string) error
}

type TodosHandler struct {
	todoService services.TodoServiceInterface
}

func NewTodoHandler(todoService services.TodoServiceInterface) *TodosHandler {
	return &TodosHandler{todoService: todoService}
}

func (h *TodosHandler) PostTodos(ctx echo.Context) error {
	var newTodo models.TodoCreate
	if err := ctx.Bind(&newTodo); err != nil {
		return pkg.SendAPIError(ctx, http.StatusBadRequest, err.Error(), err)
	}

	todo, err := h.todoService.CreateTodo(ctx.Request().Context(), newTodo)
	if err != nil {
		return pkg.SendAPIError(ctx, http.StatusBadRequest, "Invalid format for new todo", err)
	}
	return ctx.JSON(http.StatusCreated, todo)
}

func (h *TodosHandler) GetTodos(ctx echo.Context, params models.GetTodosParams) error {
	todos, err := h.todoService.GetTodos(ctx.Request().Context(), params)
	if err != nil {
		return pkg.SendAPIError(ctx, http.StatusInternalServerError, "problem getting all todos", err)
	}
	return ctx.JSON(http.StatusOK, todos)
}

func (h *TodosHandler) GetTodosId(ctx echo.Context, id string) error {
	todoUUID, err := uuid.Parse(id)
	if err != nil {
		return pkg.SendAPIError(ctx, http.StatusBadRequest, "todoId is not a valid UUID", err)
	}

	todo, err := h.todoService.GetTodoById(ctx.Request().Context(), todoUUID)
	if err != nil {
		return pkg.SendAPIError(ctx, http.StatusNotFound, "todo not found: "+id, err)
	}
	return ctx.JSON(http.StatusOK, todo)
}

func (h *TodosHandler) PutTodosId(ctx echo.Context, id string) error {
	todoUUID, err := uuid.Parse(id)
	if err != nil {
		return pkg.SendAPIError(ctx, http.StatusBadRequest, "todoId is not a valid UUID", err)
	}

	var updateTodo models.TodoUpdate
	if err := ctx.Bind(&updateTodo); err != nil {
		return pkg.SendAPIError(ctx, http.StatusBadRequest, err.Error(), err)
	}

	todo, err := h.todoService.UpdateTodo(ctx.Request().Context(), todoUUID, updateTodo)
	if err != nil {
		return pkg.SendAPIError(ctx, http.StatusInternalServerError, "failed to update todo", err)
	}
	return ctx.JSON(http.StatusOK, todo)
}

func (h *TodosHandler) DeleteTodosId(ctx echo.Context, id string) error {
	todoUUID, err := uuid.Parse(id)
	if err != nil {
		return pkg.SendAPIError(ctx, http.StatusBadRequest, "todoId is not a valid UUID", err)
	}

	err = h.todoService.DeleteTodo(ctx.Request().Context(), todoUUID)
	if err != nil {
		return pkg.SendAPIError(ctx, http.StatusInternalServerError, "failed to delete todo", err)
	}
	return ctx.NoContent(http.StatusNoContent)
}
