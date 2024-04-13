package handler

import (
	FileUC "doc-translate-go/pkg/file/usecase"
	UserENT "doc-translate-go/pkg/user/entity"
	UserUC "doc-translate-go/pkg/user/usecase"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DeleteFilesRequest struct {
	// FileIds are translated file ids.
	FileIds []int `json:"file_ids"`
}

// DeleteFiles - Delete Files
//
// @Summary Delete Files
// @Description Delete Files
// @Tags Files
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @Param file_delete_request body DeleteFilesRequest true "File Delete Request"
// @Success 200 {string} string "Files delete successfully"
// @Failure 400 {string} string "Bad request"
// @Router /delete-translated-files [delete]
func DeleteFiles(
	c echo.Context,
	userUseCase *UserUC.UserUseCase,
	originalFileMetadataUseCase *FileUC.OriginalFileMetadataUseCase,
	translatedFileUseCase *FileUC.TranslatedFileMetadataUseCase,
	fileUseCase *FileUC.FileUseCase,
) error {
	userProfileValue := c.Get("userProfile")
	user, ok := userProfileValue.(UserENT.UserProfile)
	if !ok {
		c.Logger().Error("user profile not found")
		return echo.ErrBadRequest
	}

	var req DeleteFilesRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("failed to parse request: %v", err)
		return echo.ErrBadRequest
	}

	origIds, err := translatedFileUseCase.ListOriginalFileIdsByIds(req.FileIds)
	if err != nil {
		c.Logger().Errorf("failed to get original file ids from request: %v", err)
		return echo.ErrInternalServerError
	}

	origMetas, err := originalFileMetadataUseCase.ListByIds(origIds)
	if err != nil {
		c.Logger().Errorf("failed to get original file metadata: %v", err)
		return echo.ErrInternalServerError
	}

	var origFilepaths []string
	for _, m := range origMetas {
		origFilepaths = append(origFilepaths, fmt.Sprintf("%s/%s", user.Isid, m.Filename))
	}

	err = fileUseCase.DeleteMany(origFilepaths)
	if err != nil {
		c.Logger().Errorf("failed to delete original files on S3: %v", err)
		return echo.ErrInternalServerError
	}

	translMetas, err := translatedFileUseCase.ListByIds(req.FileIds)
	if err != nil {
		c.Logger().Errorf("failed to get original file metadata: %v", err)
		return echo.ErrInternalServerError
	}

	var translFilepaths []string
	for _, m := range translMetas {
		translFilepaths = append(translFilepaths, fmt.Sprintf("%s/%s", user.Isid, m.Filename))
	}

	err = fileUseCase.DeleteMany(translFilepaths)
	if err != nil {
		c.Logger().Errorf("failed to delete translated files on S3: %v", err)
		return echo.ErrInternalServerError
	}

	err = originalFileMetadataUseCase.DeleteByIds(origIds)
	if err != nil {
		c.Logger().Errorf("failed to delete original file metadata on S3: %v", err)
		return echo.ErrInternalServerError
	}

	err = translatedFileUseCase.DeleteByIds(req.FileIds)
	if err != nil {
		c.Logger().Errorf("failed to delete translated file metadata on S3: %v", err)
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusOK, "ok")
}
