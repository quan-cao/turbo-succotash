package repository

import "doc-translate-go/pkg/file/entity"

// TranslatedFileMetadataRepository operates against a database
// or any data persistent layer.
type TranslatedFileMetadataRepository interface {
	Create(f *entity.TranslatedFileMetadata) (int, error)
	ListByIds(ids []int) ([]*entity.TranslatedFileMetadata, error)
	ListByIsid(isid string) ([]*entity.TranslatedFileMetadata, error)
	ListOriginalFileIdsByIds(ids []int) ([]int, error)
	Update(f *entity.TranslatedFileMetadata) error
	DeleteById(id int) error
	DeleteByIds(ids []int) error
}
