package postgresql

import (
	"doc-translate-go/pkg/db"
	"doc-translate-go/pkg/file/entity"
	"doc-translate-go/pkg/file/repository"
)

type PostgresqlTranslatedFileMetadataRepository struct {
	querier db.Querier
}

func NewPostgresqlTranslatedFileMetadataRepository(querier db.Querier) *PostgresqlTranslatedFileMetadataRepository {
	return &PostgresqlTranslatedFileMetadataRepository{querier}
}

func (r *PostgresqlTranslatedFileMetadataRepository) Create(f *entity.TranslatedFileMetadata) (int, error) {
	panic("")
}

func (r *PostgresqlTranslatedFileMetadataRepository) ListByIds(ids []int) ([]*entity.TranslatedFileMetadata, error) {
	panic("")
}

func (r *PostgresqlTranslatedFileMetadataRepository) ListByIsid(isid string) ([]*entity.TranslatedFileMetadata, error) {
	panic("")
}

func (r *PostgresqlTranslatedFileMetadataRepository) ListOriginalFileIdsByIds(ids []int) ([]int, error) {
	panic("")
}

func (r *PostgresqlTranslatedFileMetadataRepository) Update(f *entity.TranslatedFileMetadata) error {
	panic("")
}

func (r *PostgresqlTranslatedFileMetadataRepository) DeleteById(id int) error { panic("") }

func (r *PostgresqlTranslatedFileMetadataRepository) DeleteByIds(ids []int) error { panic("") }

// Ensure implementation
var _ repository.TranslatedFileMetadataRepository = (*PostgresqlTranslatedFileMetadataRepository)(nil)
