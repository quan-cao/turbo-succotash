package handler

import (
	"doc-translate-go/pkg/file/usecase"
	"doc-translate-go/pkg/user/entity"
	"io"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

// Translate multiple DOCX files
//
// @Summary Translate multiple DOCX files
// @Description Send multiple files to the gRPC server for translation along with source and target language selections.
// @Tags Files
// @Accept multipart/form-data
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param file formData []file true "Upload files"
// @Param sourceLang formData string true "Source Language"
// @Param targetLang formData string true "Target Language"
// @Success 200 {string} string "Files sent successfully"
// @Router /translate-docx [post]
func TranslateDocx(c echo.Context, translateUseCase *usecase.TranslateUseCase) error {
	userProfileValue := c.Get("userProfile")
	userProfile, ok := userProfileValue.(entity.UserProfile)
	if !ok {
		return echo.ErrBadRequest
	}

	form, err := c.MultipartForm()
	if err != nil {
		return echo.ErrBadRequest
	}

	files := form.File["file"]
	sourceLang := c.FormValue("sourceLang")
	targetLang := c.FormValue("targetLang")

	errChan := make(chan error)

	wg := sync.WaitGroup{}

	for _, file := range files {

		// Open files here so we raise error when encounter first error
		src, err := file.Open()
		if err != nil {
			return echo.ErrInternalServerError
		}
		defer src.Close()

		b, err := io.ReadAll(src)
		if err != nil {
			return echo.ErrInternalServerError
		}

		wg.Add(1)
		go translateFile(&wg, b, file.Filename, int(file.Size), userProfile.Isid, sourceLang, targetLang, translateUseCase, errChan)
	}

	wg.Wait()

	close(errChan)

	for e := range errChan {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": e.Error()})
	}

	return nil
}

func translateFile(
	wg *sync.WaitGroup,
	b []byte,
	filename string,
	filesize int,
	isid string,
	sourceLang string,
	targetLang string,
	translateUseCase *usecase.TranslateUseCase,
	errChan chan error,
) {
	defer wg.Done()

	err := translateUseCase.TranslateAsync(b, filename, filesize, isid, sourceLang, targetLang)
	if err != nil {
		errChan <- err
		return
	}
}
