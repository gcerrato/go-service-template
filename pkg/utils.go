package pkg

import (
	"errors"
	"log/slog"

	"github.com/gcerrato/go-service-template/api/models"
	"github.com/labstack/echo/v4"
)

func SendRepoError(msg string, err error) error {
	slog.Error(msg, "err", err.Error())
	return errors.New(msg)
}

func SendAPIError(ctx echo.Context, code int, message string, foundError error) error {
	if foundError != nil {
		slog.Error("error", "error", foundError.Error())
	}
	apiErr := models.Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, apiErr)
	return err
}
