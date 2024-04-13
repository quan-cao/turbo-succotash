package handler

import (
	_ "doc-translate-go/pkg/file/entity"
	"doc-translate-go/pkg/file/usecase"
	UserENT "doc-translate-go/pkg/user/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ShowFiles - Show all translated files user has
//
// @Summary Show all translated files
// @Description List out all translated files user has
// @Tags Files
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Success 200 {array} entity.TranslatedFileMetadata
// @Failure 400 {string} string "Bad request"
// @Router /show-translated-files [get]
func ShowTranslatedFiles(c echo.Context, translatedFileMetadataUseCase *usecase.TranslatedFileMetadataUseCase) error {
	userProfileValue := c.Get("userProfile")
	user, ok := userProfileValue.(*UserENT.UserProfile)
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
