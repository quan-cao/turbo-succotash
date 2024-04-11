package handler

import (
	"doc-translate-go/pkg/file/usecase"
	"doc-translate-go/pkg/user/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ShowTranslatedFiles(c echo.Context, translatedFileMetadataUseCase *usecase.TranslatedFileMetadataUseCase) error {
	userProfileValue := c.Get("userProfile")
	user, ok := userProfileValue.(entity.UserProfile)
	if !ok {
		c.Logger().Error("user profile not found")
		return echo.ErrBadRequest
	}

	files, err := translatedFileMetadataUseCase.ListByIsid(user.Isid)
	if err != nil {
		c.Logger().Error("failed to get translated file metadata: %v", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, files)
}
