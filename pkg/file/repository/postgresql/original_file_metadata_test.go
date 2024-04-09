package postgresql

import (
	"database/sql"
	"doc-translate-go/pkg/file/entity"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMock(t *testing.T) (sqlmock.Sqlmock, *PostgresqlOriginalFileMetadataRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresqlOriginalFileMetadataRepository(db)

	return mock, repo
}

func TestPostgresqlOriginalFileMetadataRepository__Create(t *testing.T) {
	mock, repo := newMock(t)
	ent := &entity.OriginalFileMetadata{Id: 1}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(ent.Id)

	mock.ExpectQuery("INSERT INTO original_files").WillReturnRows(rows)

	got, err := repo.Create(ent)
	if err != nil {
		t.Fatal(err)
	}

	if got != ent.Id {
		t.Fatalf("expect %v, got %v", ent.Id, got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__Create_ScanErr(t *testing.T) {
	mock, repo := newMock(t)
	ent := &entity.OriginalFileMetadata{Id: 1}

	e := sql.ErrNoRows
	rows := sqlmock.NewRows([]string{"id"})

	mock.ExpectQuery("INSERT INTO original_files").WillReturnRows(rows).WillReturnError(e)

	got, err := repo.Create(ent)
	if err != e {
		t.Fatalf("expect %v, got %v", e, err)
	}

	if got != 0 {
		t.Fatalf("expect %v, got %v", 0, got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__ListByIds(t *testing.T) {
	mock, repo := newMock(t)

	want := []*entity.OriginalFileMetadata{{Id: 1}, {Id: 2}}

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)
	for _, r := range want {
		rows.AddRow(r.Id, r.SHA256, r.Filename, r.FileType, r.FileSize, r.SourceLanguage, r.TokenCount, r.CreatedAt, r.UpdatedAt, r.CreatedBy)
	}

	mock.ExpectQuery("SELECT .+ FROM original_files WHERE id IN \\([$0-9, ]+\\);").WillReturnRows(rows)

	got, err := repo.ListByIds([]int{1, 2})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__ListByIds_Err(t *testing.T) {
	mock, repo := newMock(t)

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	e := errors.New("error list by ids")
	mock.ExpectQuery("SELECT .+ FROM original_files WHERE id IN \\([$0-9, ]+\\);").
		WillReturnRows(rows).
		WillReturnError(e)

	got, err := repo.ListByIds([]int{1, 2})
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__ListByIsid(t *testing.T) {
	mock, repo := newMock(t)

	want := []*entity.OriginalFileMetadata{{Id: 1}, {Id: 2}}

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	for _, i := range want {
		rows.AddRow(i.Id, i.SHA256, i.Filename, i.FileType, i.FileSize, i.SourceLanguage, i.TokenCount, i.CreatedAt, i.UpdatedAt, i.CreatedBy)
	}

	mock.ExpectQuery("SELECT .+ FROM original_files WHERE created_by = \\$1;").WillReturnRows(rows)

	got, err := repo.ListByIsid("1")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("exptected %v, got %v", want, got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__ListByIsid_Err(t *testing.T) {
	mock, repo := newMock(t)

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	e := errors.New("list error")
	mock.ExpectQuery("SELECT .+ FROM original_files WHERE created_by = \\$1;").WillReturnRows(rows).WillReturnError(e)

	got, err := repo.ListByIsid("1")
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != nil {
		t.Fatalf("exptected nil, got %v", got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__ListByFilenameIsid(t *testing.T) {
	mock, repo := newMock(t)

	want := []*entity.OriginalFileMetadata{{Id: 1}, {Id: 2}}

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	for _, i := range want {
		rows.AddRow(i.Id, i.SHA256, i.Filename, i.FileType, i.FileSize, i.SourceLanguage, i.TokenCount, i.CreatedAt, i.UpdatedAt, i.CreatedBy)
	}

	mock.ExpectQuery("SELECT .+ FROM original_files WHERE file_name = \\$1 AND created_by = \\$2;").WillReturnRows(rows)

	got, err := repo.ListByFilenameIsid("file", "1")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("exptected %v, got %v", want, got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__ListByFilenameIsid_Err(t *testing.T) {
	mock, repo := newMock(t)

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)

	e := errors.New("list error")
	mock.ExpectQuery("SELECT .+ FROM original_files WHERE file_name = \\$1 AND created_by = \\$2;").WillReturnRows(rows).WillReturnError(e)

	got, err := repo.ListByFilenameIsid("file", "1")
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != nil {
		t.Fatalf("exptected nil, got %v", got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__Update(t *testing.T) {
	mock, repo := newMock(t)
	ent := &entity.OriginalFileMetadata{}

	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec("UPDATE original_files SET .+ WHERE id = \\$\\d{1,2};").WillReturnResult(result)

	err := repo.Update(ent)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__Update_Err(t *testing.T) {
	mock, repo := newMock(t)
	ent := &entity.OriginalFileMetadata{}

	result := sqlmock.NewResult(1, 1)
	e := errors.New("update error")
	mock.ExpectExec("UPDATE original_files SET .+ WHERE id = \\$\\d{1,2};").WillReturnResult(result).WillReturnError(e)

	err := repo.Update(ent)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__DeleteById(t *testing.T) {
	mock, repo := newMock(t)
	id := 1

	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec("DELETE FROM original_files WHERE id = \\$1;").WillReturnResult(result)

	err := repo.DeleteById(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__DeleteById_Err(t *testing.T) {
	mock, repo := newMock(t)
	id := 1

	result := sqlmock.NewResult(1, 1)
	e := errors.New("delete error")
	mock.ExpectExec("DELETE FROM original_files WHERE id = \\$1;").WillReturnResult(result).WillReturnError(e)

	err := repo.DeleteById(id)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__DeleteByIds(t *testing.T) {
	mock, repo := newMock(t)
	ids := []int{1, 2, 3}

	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec("DELETE FROM original_files WHERE id IN \\([$\\d, ]+\\);").WillReturnResult(result)

	err := repo.DeleteByIds(ids)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__DeleteByIds_Err(t *testing.T) {
	mock, repo := newMock(t)
	ids := []int{1, 2, 3}

	result := sqlmock.NewResult(1, 1)
	e := errors.New("delete error")
	mock.ExpectExec("DELETE FROM original_files WHERE id IN \\([$\\d, ]+\\);").WillReturnResult(result).WillReturnError(e)

	err := repo.DeleteByIds(ids)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}
