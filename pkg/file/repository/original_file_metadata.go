package repository

import "doc-translate-go/pkg/file/entity"

// OriginalFileMetadataRepository operates against a database
// or any data persistent layer.
type OriginalFileMetadataRepository interface {
	Create(f *entity.OriginalFileMetadata) (int, error)
	ListByIds(ids []int) ([]*entity.OriginalFileMetadata, error)
	ListByIsid(isid string) ([]*entity.OriginalFileMetadata, error)
	ListByFilenameIsid(filename string, isid string) ([]*entity.OriginalFileMetadata, error)
	Update(f *entity.OriginalFileMetadata) error
	DeleteById(id int) error
	DeleteByIds(ids []int) error
}
