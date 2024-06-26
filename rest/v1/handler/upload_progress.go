package handler

import (
	"doc-translate-go/pkg/file/usecase"
	"doc-translate-go/pkg/user/entity"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Status struct {
	File       string    `json:"file"`
	Status     string    `json:"status"`
	SourceLang string    `json:"source_lang"`
	TargetLang string    `json:"target_lang"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UploadProgressResponse struct {
	Isid     string    `json:"isid"`
	Statuses []*Status `json:"files_status"`
}

// @Summary Get File upload progress
// @Description Retrieve the file upload progress based on isid
// @Tags Files
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string ture "Authorization"
// @Success 200 {string} string UploadProgressResponse
// @Router /upload-progress [get]
func UploadProgress(c echo.Context, progressUseCase *usecase.ProgressUseCase) error {
	userProfileValue := c.Get("userProfile")
	user, ok := userProfileValue.(*entity.UserProfile)
	if !ok {
		c.Logger().Error("user profile not found")
		return echo.ErrBadRequest
	}

	statuses, err := progressUseCase.ListByIsid(user.Isid)
	if err != nil {
		c.Logger().Errorf("unable to get file status: %v", err)
		return echo.ErrInternalServerError
	}

	resp := UploadProgressResponse{Isid: user.Isid}

	for _, s := range statuses {
		stt := &Status{
			File:       strings.Replace(s.Key, fmt.Sprintf("%s_", user.Isid), "", 1),
			Status:     strings.Replace(s.Status, "fail:", "", 1),
			UpdatedAt:  s.UpdatedAt,
			SourceLang: s.SourceLang,
			TargetLang: s.TargetLang,
		}

		resp.Statuses = append(resp.Statuses, stt)
	}

	return c.JSON(http.StatusOK, resp)
}
