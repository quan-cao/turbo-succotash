package usecase

import (
	"doc-translate-go/pkg/tracker"
	"fmt"
)

type ProgressUseCase struct {
	fileTracker tracker.FileTracker
}

func NewProgressUseCase(tracker tracker.FileTracker) *ProgressUseCase {
	return &ProgressUseCase{tracker}
}

func (uc *ProgressUseCase) ListByIsid(isid string) ([]*tracker.FileStatus, error) {
	return uc.fileTracker.List(fmt.Sprintf("%s*", isid))
}
