package usecase

import (
	"doc-translate-go/pkg/file/repository"
)

type FileUseCase struct {
	repo repository.FileRepository
}

func NewFileUseCase(repo repository.FileRepository) *FileUseCase {
	return &FileUseCase{repo}
}

func (uc *FileUseCase) Get(filepath string) ([]byte, error) {
	return uc.repo.Get(filepath)
}

func (uc *FileUseCase) Persist(b []byte, filepath string) error {
	return uc.repo.Persist(b, filepath)
}

func (uc *FileUseCase) Delete(filepath string) error {
	return uc.repo.Delete(filepath)
}
