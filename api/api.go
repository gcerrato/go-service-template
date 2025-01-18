//go:generate go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=models.config.yaml api.yaml
//go:generate go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=server.config.yaml api.yaml

package api

import (
	_ "embed"

	"github.com/gcerrato/go-service-template/internal/handlers"
	"github.com/gcerrato/go-service-template/internal/services"
)

//go:embed api.yaml
var Spec []byte

type ServerHandler struct {
	handlers.TodosAPI
}

func NewServerHandler(todoService services.TodoService) *ServerHandler {
	return &ServerHandler{
		TodosAPI: handlers.NewTodoHandler(&todoService),
	}
}

func (h *ServerHandler) GetSwagger() []byte {
	return Spec
}
