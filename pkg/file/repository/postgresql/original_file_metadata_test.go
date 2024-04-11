package postgresql

import (
	"database/sql"
	"doc-translate-go/pkg/file/entity"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newOrigMock(t *testing.T) (sqlmock.Sqlmock, *PostgresqlOriginalFileMetadataRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresqlOriginalFileMetadataRepository(db)

	return mock, repo
}

func TestPostgresqlOriginalFileMetadataRepository_Create(t *testing.T) {
	mock, repo := newOrigMock(t)
	ent := &entity.OriginalFileMetadata{Id: 1}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(ent.Id)

	cmd := `INSERT INTO original_files \(sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by\)
                VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9\)
                RETURNING id;`

	mock.ExpectQuery(cmd).WillReturnRows(rows)

	got, err := repo.Create(ent)
	if err != nil {
		t.Fatal(err)
	}

	if got != ent.Id {
		t.Fatalf("expected %v, got %v", ent.Id, got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_Create_ScanErr(t *testing.T) {
	mock, repo := newOrigMock(t)
	ent := &entity.OriginalFileMetadata{Id: 1}

	e := sql.ErrNoRows
	rows := sqlmock.NewRows([]string{"id"})

	cmd := `INSERT INTO original_files \(sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by\)
                VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9\)
                RETURNING id;`

	mock.ExpectQuery(cmd).WillReturnRows(rows).WillReturnError(e)

	got, err := repo.Create(ent)
	if err != e {
		t.Fatalf("expect %v, got %v", e, err)
	}

	if got != 0 {
		t.Fatalf("expect %v, got %v", 0, got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_ListByIds(t *testing.T) {
	mock, repo := newOrigMock(t)

	want := []*entity.OriginalFileMetadata{{Id: 1}, {Id: 2}}

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)
	for _, r := range want {
		rows.AddRow(r.Id, r.SHA256, r.Filename, r.FileType, r.FileSize, r.SourceLanguage, r.TokenCount, r.CreatedAt, r.UpdatedAt, r.CreatedBy)
	}

	cmd := `SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
                FROM original_files
                WHERE id IN \([$0-9, ]+\);`

	mock.ExpectQuery(cmd).WillReturnRows(rows)

	got, err := repo.ListByIds([]int{1, 2})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_ListByIds_Err(t *testing.T) {
	mock, repo := newOrigMock(t)

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	cmd := `SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
                FROM original_files
                WHERE id IN \([$0-9, ]+\);`

	e := errors.New("error list by ids")
	mock.ExpectQuery(cmd).
		WillReturnRows(rows).
		WillReturnError(e)

	got, err := repo.ListByIds([]int{1, 2})
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_ListByIsid(t *testing.T) {
	mock, repo := newOrigMock(t)

	want := []*entity.OriginalFileMetadata{{Id: 1}, {Id: 2}}

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	for _, i := range want {
		rows.AddRow(i.Id, i.SHA256, i.Filename, i.FileType, i.FileSize, i.SourceLanguage, i.TokenCount, i.CreatedAt, i.UpdatedAt, i.CreatedBy)
	}

	cmd := `SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
        FROM original_files
        WHERE created_by = \$1;`

	mock.ExpectQuery(cmd).WillReturnRows(rows)

	got, err := repo.ListByIsid("1")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("exptected %v, got %v", want, got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_ListByIsid_Err(t *testing.T) {
	mock, repo := newOrigMock(t)

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	cmd := `SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
        FROM original_files
        WHERE created_by = \$1;`

	e := errors.New("list error")
	mock.ExpectQuery(cmd).WillReturnRows(rows).WillReturnError(e)

	got, err := repo.ListByIsid("1")
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != nil {
		t.Fatalf("exptected nil, got %v", got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_ListByFilenameIsid(t *testing.T) {
	mock, repo := newOrigMock(t)

	want := []*entity.OriginalFileMetadata{{Id: 1}, {Id: 2}}

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	for _, i := range want {
		rows.AddRow(i.Id, i.SHA256, i.Filename, i.FileType, i.FileSize, i.SourceLanguage, i.TokenCount, i.CreatedAt, i.UpdatedAt, i.CreatedBy)
	}

	cmd := `SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
        FROM original_files
        WHERE file_name = \$1 AND created_by = \$2;`

	mock.ExpectQuery(cmd).WillReturnRows(rows)

	got, err := repo.ListByFilenameIsid("file", "1")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("exptected %v, got %v", want, got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_ListByFilenameIsid_Err(t *testing.T) {
	mock, repo := newOrigMock(t)

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	cmd := `SELECT id, sha256, filename, file_type, file_size, source_language, token_count, created_at, updated_at, created_by
        FROM original_files
        WHERE file_name = \$1 AND created_by = \$2;`

	e := errors.New("list error")
	mock.ExpectQuery(cmd).WillReturnRows(rows).WillReturnError(e)

	got, err := repo.ListByFilenameIsid("file", "1")
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != nil {
		t.Fatalf("exptected nil, got %v", got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_Update(t *testing.T) {
	mock, repo := newOrigMock(t)
	ent := &entity.OriginalFileMetadata{}

	cmd := `UPDATE original_files SET
                sha256 = \$1,
                filename = \$2,
                file_type = \$3,
                file_size = \$4,
                source_language = \$5,
                token_count = \$6,
                updated_at = \$7,
                created_by = \$8
        WHERE id = \$9;`

	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(cmd).WillReturnResult(result)

	err := repo.Update(ent)
	if err != nil {
		t.Fatal(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_Update_Err(t *testing.T) {
	mock, repo := newOrigMock(t)
	ent := &entity.OriginalFileMetadata{}

	cmd := `UPDATE original_files SET
                sha256 = \$1,
                filename = \$2,
                file_type = \$3,
                file_size = \$4,
                source_language = \$5,
                token_count = \$6,
                updated_at = \$7,
                created_by = \$8
        WHERE id = \$9;`

	result := sqlmock.NewResult(1, 1)
	e := errors.New("update error")
	mock.ExpectExec(cmd).WillReturnResult(result).WillReturnError(e)

	err := repo.Update(ent)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_DeleteById(t *testing.T) {
	mock, repo := newOrigMock(t)
	id := 1

	cmd := `DELETE FROM original_files WHERE id = \$1;`

	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(cmd).WillReturnResult(result)

	err := repo.DeleteById(id)
	if err != nil {
		t.Fatal(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_DeleteById_Err(t *testing.T) {
	mock, repo := newOrigMock(t)
	id := 1

	cmd := `DELETE FROM original_files WHERE id = \$1;`

	result := sqlmock.NewResult(1, 1)
	e := errors.New("delete error")
	mock.ExpectExec(cmd).WillReturnResult(result).WillReturnError(e)

	err := repo.DeleteById(id)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_DeleteByIds(t *testing.T) {
	mock, repo := newOrigMock(t)
	ids := []int{1, 2, 3}

	cmd := `DELETE FROM original_files WHERE id IN \([$\d, ]+\);`

	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(cmd).WillReturnResult(result)

	err := repo.DeleteByIds(ids)
	if err != nil {
		t.Fatal(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlOriginalFileMetadataRepository_DeleteByIds_Err(t *testing.T) {
	mock, repo := newOrigMock(t)
	ids := []int{1, 2, 3}

	cmd := `DELETE FROM original_files WHERE id IN \([$\d, ]+\);`

	result := sqlmock.NewResult(1, 1)
	e := errors.New("delete error")
	mock.ExpectExec(cmd).WillReturnResult(result).WillReturnError(e)

	err := repo.DeleteByIds(ids)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}
