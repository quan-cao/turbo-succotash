package handler

import (
	"doc-translate-go/pkg/user/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Authorize godoc
//
// @Summary Authorize user
// @Description Initiates user authorization process by redirecting to the authorization service.
// @Tags User
// @Accept json
// @Produce json
// @Success 302 {header} string Location "Redirect to authorization URL"
// @Failure 400 {object} map[string]any "Bad Request: Invalid input parameters"
// @Failure 500 {object} map[string]any "Internal Server Error: Processing error"
// @Router /authorize [get]
func Authorize(c echo.Context, authUseCase *usecase.AuthUseCase) error {
	redirectUrl, err := authUseCase.Authorize()
	if err != nil {
		c.Logger().Errorf("failed to authorize: %v", err)
		return echo.ErrBadRequest
	}
	return c.Redirect(http.StatusFound, redirectUrl)
}
