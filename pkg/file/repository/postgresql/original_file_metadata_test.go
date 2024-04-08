package postgresql

import (
	"database/sql"
	"doc-translate-go/pkg/file/entity"
	"reflect"
	"testing"
	"time"

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

	mock.ExpectQuery("INSERT INTO original_files").WillReturnRows(rows)

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

	want := []*entity.OriginalFileMetadata{
		{
			Id:             1,
			SHA256:         "",
			Filename:       "",
			FileType:       "docx",
			FileSize:       0,
			SourceLanguage: "",
			TokenCount:     0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			CreatedBy:      "1",
		},
	}

	columns := []string{"id", "sha256", "filename", "file_type", "file_size", "source_language", "token_count", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)
	for _, r := range want {
		rows.AddRow(r.Id, r.SHA256, r.Filename, r.FileType, r.FileSize, r.SourceLanguage, r.TokenCount, r.CreatedAt, r.UpdatedAt, r.CreatedBy)
	}

	mock.ExpectQuery("SELECT .+ FROM original_files WHERE id IN").WillReturnRows(rows)

	got, err := repo.ListByIds([]int{1})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
