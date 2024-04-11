package handler

import (
	"doc-translate-go/pkg/user/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Authorize(c echo.Context, userUseCase *usecase.UserUseCase) error {
	redirectUrl, err := usecase.Authorize()
	if err != nil {
		c.Logger().Errorf("failed to authorize: %v", err)
		return echo.ErrBadRequest
	}
	return c.Redirect(http.StatusFound, redirectUrl)
}
