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

func (uc *TranslatedFileMetadataUseCase) ListByIds(ids []int) ([]*entity.TranslatedFileMetadata, error) {
	return uc.repo.ListByIds(ids)
}

func (uc *TranslatedFileMetadataUseCase) ListOriginalFileIdsByIds(ids []int) ([]int, error) {
	return uc.repo.ListOriginalFileIdsByIds(ids)
}

func (uc *TranslatedFileMetadataUseCase) DeleteByIds(ids []int) error {
	return uc.repo.DeleteByIds(ids)
}
