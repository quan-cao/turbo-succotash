package usecase

import (
	"doc-translate-go/pkg/file/entity"
	"doc-translate-go/pkg/file/repository"
)

type TranslatedFileMetadataUseCase struct {
	repo repository.TranslatedFileMetadataRepository
}

func NewTranslatedFileMetadataUseCase(repo repository.TranslatedFileMetadataRepository) *TranslatedFileMetadataUseCase {
	return &TranslatedFileMetadataUseCase{repo}
}

func (uc *TranslatedFileMetadataUseCase) Persist(f *entity.TranslatedFileMetadata) (int, error) {
	return uc.repo.Create(f)
}

func (uc *TranslatedFileMetadataUseCase) ListByIsid(isid string) ([]*entity.TranslatedFileMetadata, error) {
	return uc.repo.ListByIsid(isid)
}
