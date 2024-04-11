package postgresql

import (
	"doc-translate-go/pkg/db"
	"doc-translate-go/pkg/file/entity"
	"doc-translate-go/pkg/file/repository"
	"fmt"
	"strings"
)

type PostgresqlTranslatedFileMetadataRepository struct {
	querier db.Querier
}

func NewPostgresqlTranslatedFileMetadataRepository(querier db.Querier) *PostgresqlTranslatedFileMetadataRepository {
	return &PostgresqlTranslatedFileMetadataRepository{querier}
}

func (r *PostgresqlTranslatedFileMetadataRepository) Create(f *entity.TranslatedFileMetadata) (int, error) {
	cmd := `INSERT INTO translated_files (original_files_id, translated_filename, target_language, cost, time_taken, created_at, updated_at, created_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id;`

	row := r.querier.QueryRow(cmd, f.OriginalFileId, f.Filename, f.TargetLanguage, f.Cost, f.TimeTaken, f.CreatedAt, f.UpdatedAt, f.CreatedBy)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresqlTranslatedFileMetadataRepository) ListByIds(ids []int) ([]*entity.TranslatedFileMetadata, error) {
	var out []*entity.TranslatedFileMetadata

	arg_placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		arg_placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	cmd := fmt.Sprintf(`
                SELECT id, original_files_id, translated_filename, target_language, cost, time_taken, created_at, updated_at, created_by
                FROM translated_files
                WHERE id IN (%s);`,
		strings.Join(arg_placeholders, ", "),
	)

	rows, err := r.querier.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f entity.TranslatedFileMetadata
		if err := rows.Scan(&f.Id, &f.OriginalFileId, &f.Filename, &f.TargetLanguage, &f.Cost, &f.TimeTaken, &f.CreatedAt, &f.UpdatedAt, &f.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PostgresqlTranslatedFileMetadataRepository) ListByIsid(isid string) ([]*entity.TranslatedFileMetadata, error) {
	var out []*entity.TranslatedFileMetadata

	cmd := `SELECT id, original_files_id, translated_filename, target_language, cost, time_taken, created_by, updated_at, created_by
                FROM translated_files
                WHERE created_by = $1;`

	rows, err := r.querier.Query(cmd, isid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f entity.TranslatedFileMetadata
		if err := rows.Scan(&f.Id, &f.OriginalFileId, &f.Filename, &f.TargetLanguage, &f.Cost, &f.TimeTaken, &f.CreatedAt, &f.UpdatedAt, &f.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PostgresqlTranslatedFileMetadataRepository) ListOriginalFileIdsByIds(ids []int) ([]int, error) {
	var out []int

	arg_placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		arg_placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	cmd := fmt.Sprintf(`
                SELECT original_files_id FROM translated_files WHERE id IN (%s);`,
		strings.Join(arg_placeholders, ", "),
	)

	rows, err := r.querier.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PostgresqlTranslatedFileMetadataRepository) Update(f *entity.TranslatedFileMetadata) error {
	cmd := `UPDATE translated_files SET
                original_files_id = $1,
                translated_filename = $2,
                target_language = $3,
                cost = $4,
                time_taken = $5,
                created_at = $6,
                updated_at = $7,
                created_by = $8
        WHERE id = $9;`

	_, err := r.querier.Exec(cmd, f.OriginalFileId, f.Filename, f.TargetLanguage, f.Cost, f.TimeTaken, f.CreatedAt, f.UpdatedAt, f.CreatedBy, f.Id)

	return err
}

func (r *PostgresqlTranslatedFileMetadataRepository) DeleteById(id int) error {
	cmd := `DELETE FROM translated_files WHERE id = $1;`
	_, err := r.querier.Exec(cmd, id)
	return err
}

func (r *PostgresqlTranslatedFileMetadataRepository) DeleteByIds(ids []int) error {
	arg_placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		arg_placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	cmd := fmt.Sprintf(
		`DELETE FROM translated_files WHERE id IN (%s);`,
		strings.Join(arg_placeholders, ", "),
	)

	_, err := r.querier.Exec(cmd, args...)
	return err
}

// Ensure implementation
var _ repository.TranslatedFileMetadataRepository = (*PostgresqlTranslatedFileMetadataRepository)(nil)
