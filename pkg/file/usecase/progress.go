package usecase

import (
	"doc-translate-go/pkg/tracker"
	"fmt"
)

type ProgressUseCase struct {
	fileTracker tracker.FileTracker
}

func (uc *ProgressUseCase) ListByIsid(isid string) ([]*tracker.FileStatus, error) {
	return uc.fileTracker.List(fmt.Sprintf("%s*", isid))
}
