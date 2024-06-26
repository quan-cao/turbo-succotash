package handler

import (
	"doc-translate-go/pkg/user/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TokenRequest struct {
	GrantType string `json:"grant_type"`
	Token     string `json:"token"`
}

// @Sumaary Retrieve Access Token
// @Description Retrieves access token based on the provided grant type and token.
// @Tags User
// @Accept json
// @Produce json
// @Param request body TokenRequest true "Token Request"
// @Success 200 {object} map[string]any "Access Token Response"
// @Router /token [post]
func Token(c echo.Context, authUseCase *usecase.AuthUseCase) error {
	var req TokenRequest

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("failed to parse request: %v", err)
		return echo.ErrBadRequest
	}

	accessToken, err := authUseCase.RetrieveAccessToken(req.GrantType, req.Token)
	if err != nil {
		c.Logger().Errorf("failed to retrieve access token: %v", err)
		return echo.ErrBadRequest
	}

	user, err := authUseCase.RetrieveUserProfile(accessToken)
	if err != nil {
		c.Logger().Errorf("failed to retrieve user profile: %v", err)
		return echo.ErrInternalServerError
	}

	err = authUseCase.ValidateDistributionListHasIsid(user.Isid)
	if err != nil {
		c.Logger().Errorf("user not in distribution list: %v", err)
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, map[string]any{"access_token": accessToken})
}
