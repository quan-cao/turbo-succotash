package postgresql

import (
	"doc-translate-go/pkg/db"
	"doc-translate-go/pkg/file/entity"
	"doc-translate-go/pkg/file/repository"
	"fmt"
	"strings"
)

type PostgresqlOriginalFileMetadataRepository struct {
	querier db.Querier
}

func NewPostgresqlOriginalFileMetadataRepository(querier db.Querier) *PostgresqlOriginalFileMetadataRepository {
	return &PostgresqlOriginalFileMetadataRepository{querier}
}

func (r *PostgresqlOriginalFileMetadataRepository) Create(f *entity.OriginalFileMetadata) (int, error) {
	cmd := `INSERT INTO original_files (sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                RETURNING id;`

	row := r.querier.QueryRow(cmd, f.SHA256, f.Filename, f.FileType, f.FileSize, f.SourceLanguage, f.TokenCount, f.CreatedAt, f.UpdatedAt, f.CreatedBy)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresqlOriginalFileMetadataRepository) ListByIds(ids []int) ([]*entity.OriginalFileMetadata, error) {
	var out []*entity.OriginalFileMetadata

	arg_placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		arg_placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	cmd := fmt.Sprintf(`
                SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
                FROM original_files
                WHERE id IN (%s);`,
		strings.Join(arg_placeholders, ", "),
	)

	rows, err := r.querier.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f entity.OriginalFileMetadata
		if err := rows.Scan(&f.Id, &f.SHA256, &f.Filename, &f.FileType, &f.FileSize, &f.SourceLanguage, &f.TokenCount, &f.CreatedAt, &f.UpdatedAt, &f.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PostgresqlOriginalFileMetadataRepository) ListByIsid(isid string) ([]*entity.OriginalFileMetadata, error) {
	var out []*entity.OriginalFileMetadata

	cmd := `SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
        FROM original_files
        WHERE created_by = $1;`

	rows, err := r.querier.Query(cmd, isid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f entity.OriginalFileMetadata
		if err := rows.Scan(&f.Id, &f.SHA256, &f.Filename, &f.FileType, &f.FileSize, &f.SourceLanguage, &f.TokenCount, &f.CreatedAt, &f.UpdatedAt, &f.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PostgresqlOriginalFileMetadataRepository) ListByFilenameIsid(filename string, isid string) ([]*entity.OriginalFileMetadata, error) {
	var out []*entity.OriginalFileMetadata

	cmd := `SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
        FROM original_files
        WHERE file_name = $1 AND created_by = $2;`

	rows, err := r.querier.Query(cmd, filename, isid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f entity.OriginalFileMetadata
		if err := rows.Scan(&f.Id, &f.SHA256, &f.Filename, &f.FileType, &f.FileSize, &f.SourceLanguage, &f.TokenCount, &f.CreatedAt, &f.UpdatedAt, &f.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PostgresqlOriginalFileMetadataRepository) Update(f *entity.OriginalFileMetadata) error {
	cmd := `UPDATE original_file SET
                sha256 = $1,
                filename = $2,
                file_type = $3,
                file_size = $4,
                source_language = $5,
                token_count = $6,
                updated_at = $7,
                created_by = $8
        WHERE id = $9;`

	_, err := r.querier.Exec(cmd, f.SHA256, f.Filename, f.FileType, f.FileSize, f.SourceLanguage, f.TokenCount, f.UpdatedAt, f.CreatedBy)
	return err
}

func (r *PostgresqlOriginalFileMetadataRepository) DeleteById(id int) error {
	cmd := `DELETE FROM original_file WHERE id = $1;`

	_, err := r.querier.Exec(cmd, id)
	return err
}

func (r *PostgresqlOriginalFileMetadataRepository) DeleteByIds(ids []int) error {
	arg_placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		arg_placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	cmd := fmt.Sprintf(`DELETE FROM original_files WHERE id IN (%s)`, arg_placeholders)

	_, err := r.querier.Exec(cmd, args...)
	return err
}

// Ensure implementation
var _ repository.OriginalFileMetadataRepository = (*PostgresqlOriginalFileMetadataRepository)(nil)
