package handler

import (
	"archive/zip"
	"bytes"
	"doc-translate-go/pkg/file/usecase"
	"doc-translate-go/pkg/user/entity"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DownloadTranslatedFiles(c echo.Context, translatedFileUseCase *usecase.TranslatedFileMetadataUseCase, fileUseCase *usecase.FileUseCase) error {
	type Request struct {
		FileIds []int `json:"file_ids"`
	}

	type Response struct {
		ZipData []byte `json:"zipData"`
	}

	userProfileValue := c.Get("userProfile")
	user, ok := userProfileValue.(entity.UserProfile)
	if !ok {
		c.Logger().Error("user profile not found")
		return echo.ErrBadRequest
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("failed to parse request: %v", err)
		return echo.ErrInternalServerError
	}

	translMetas, err := translatedFileUseCase.ListByIds(req.FileIds)
	if err != nil {
		c.Logger().Error("failed to get translated file metadata: %v", err)
		return echo.ErrInternalServerError
	}

	var filepaths []string
	for _, m := range translMetas {
		filepaths = append(filepaths, fmt.Sprintf("%s/%s", user.Isid, m.Filename))
	}

	data, errors := fileUseCase.GetMany(filepaths)
	if len(errors) > 0 {
		c.Logger().Errorf("failed to get files: %v", errors)
		return echo.ErrInternalServerError
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for p, dat := range data {
		f, err := zipWriter.Create(p)
		if err != nil {
			c.Logger().Errorf("failed to create zip writer for file %v: %v", p, errors)
			return echo.ErrInternalServerError
		}

		_, err = f.Write(dat)
		if err != nil {
			c.Logger().Errorf("failed to create write zip for file %v: %v", p, errors)
			return echo.ErrInternalServerError
		}
	}

	if err := zipWriter.Close(); err != nil {
		c.Logger().Errorf("failed to finalize zip: %v", errors)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &Response{ZipData: buf.Bytes()})
}
