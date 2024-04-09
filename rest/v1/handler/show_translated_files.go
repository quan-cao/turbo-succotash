package handler

import (
	"doc-translate-go/pkg/file/usecase"
	"doc-translate-go/pkg/user/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ShowTranslatedFiles(c echo.Context, translatedFileUseCase *usecase.TranslatedFileMetadataUseCase) error {
	userProfileValue := c.Get("userProfile")
	userProfile, ok := userProfileValue.(entity.UserProfile)
	if !ok {
		return echo.ErrBadRequest
	}

	files, err := translatedFileUseCase.ListByIsid(userProfile.Isid)
	if err != nil {
		return echo.ErrNotFound
	}

	c.JSON(http.StatusOK, files)
	return nil
}
