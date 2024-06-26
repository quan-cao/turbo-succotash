package middleware

import (
	"doc-translate-go/pkg/user/entity"
	"doc-translate-go/pkg/user/usecase"
	"os"

	"github.com/labstack/echo/v4"
)

var devToken = os.Getenv("DEV_TOKEN")

func AuthMiddleware(userUseCase *usecase.UserUseCase, authUseCase *usecase.AuthUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")

			var userProfile *entity.UserProfile
			var err error

			if token == devToken {
				userProfile = &entity.UserProfile{
					Isid:       "developer",
					GivenName:  "developer",
					FamilyName: "developer",
					Email:      "developer@merck.com",
				}
			} else {
				userProfile, err = authUseCase.RetrieveUserProfile(token)
				if err != nil {
					c.Logger().Errorf("unable to retrieve user profile: %v", err)
					return echo.ErrUnauthorized
				}
			}

			user, err := userUseCase.GetByIsid(userProfile.Isid)
			if err != nil || user == nil {
				newUser := entity.User{
					Isid:  userProfile.Isid,
					Email: userProfile.Email,
					Role:  "user",
				}

				_, err := userUseCase.Persist(&newUser)
				if err != nil {
					c.Logger().Errorf("unable to persist new user: %v", err)
					return echo.ErrInternalServerError
				}
			}

			c.Set("userProfile", userProfile)
			return next(c)
		}
	}
}
