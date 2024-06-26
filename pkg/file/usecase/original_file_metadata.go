package usecase

import (
	"doc-translate-go/pkg/file/entity"
	"doc-translate-go/pkg/file/repository"
)

type OriginalFileMetadataUseCase struct {
	repo repository.OriginalFileMetadataRepository
}

func NewOriginalFileMetadataUseCase(repo repository.OriginalFileMetadataRepository) *OriginalFileMetadataUseCase {
	return &OriginalFileMetadataUseCase{repo}
}

func (uc *OriginalFileMetadataUseCase) Persist(f *entity.OriginalFileMetadata) (int, error) {
	return uc.repo.Create(f)
}

func (uc *OriginalFileMetadataUseCase) ListByIds(ids []int) ([]*entity.OriginalFileMetadata, error) {
	return uc.repo.ListByIds(ids)
}

func (uc *OriginalFileMetadataUseCase) ListByFilenameIsid(filename string, isid string) ([]*entity.OriginalFileMetadata, error) {
	return uc.repo.ListByFilenameIsid(filename, isid)
}

func (uc *OriginalFileMetadataUseCase) DeleteByIds(ids []int) error {
	return uc.repo.DeleteByIds(ids)
}
