package postgresql

import (
	"doc-translate-go/pkg/file/entity"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTranslMock(t *testing.T) (sqlmock.Sqlmock, *PostgresqlTranslatedFileMetadataRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresqlTranslatedFileMetadataRepository(db)

	return mock, repo
}

func TestPostgresqlTranslatedFileMetadataRepository_Create(t *testing.T) {
	mock, repo := newTranslMock(t)
	ent := &entity.TranslatedFileMetadata{Id: 1}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(ent.Id)

	cmd := `INSERT INTO translated_files \(original_files_id, translated_filename, target_language, cost, time_taken, created_at, updated_at, created_by\)
        VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)
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

func TestPostgresqlTranslatedFileMetadataRepository_Create_Err(t *testing.T) {
	mock, repo := newTranslMock(t)
	ent := &entity.TranslatedFileMetadata{}

	cmd := `INSERT INTO translated_files \(original_files_id, translated_filename, target_language, cost, time_taken, created_at, updated_at, created_by\)
        VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)
        RETURNING id;`

	e := errors.New("create err")
	mock.ExpectQuery(cmd).WillReturnRows().WillReturnError(e)

	got, err := repo.Create(ent)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != 0 {
		t.Fatalf("expected %v, got %v", 0, got)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlTranslatedFileMetadataRepository_ListByIds(t *testing.T) {
	mock, repo := newTranslMock(t)

	want := []*entity.TranslatedFileMetadata{{Id: 1}, {Id: 2}}

	columns := []string{"id", "original_files_id", "translated_filename", "target_language", "cost", "time_taken", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)
	for _, r := range want {
		rows.AddRow(r.Id, r.OriginalFileId, r.Filename, r.TargetLanguage, r.Cost, r.TimeTaken, r.CreatedAt, r.UpdatedAt, r.CreatedBy)
	}

	cmd := `SELECT id, original_files_id, translated_filename, target_language, cost, time_taken, created_at, updated_at, created_by
                FROM translated_files
                WHERE id IN \([$\d, ]+\);`

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

func TestPostgresqlTranslatedFileMetadataRepository_ListByIds_Err(t *testing.T) {
	mock, repo := newTranslMock(t)

	cmd := `SELECT id, original_files_id, translated_filename, target_language, cost, time_taken, created_at, updated_at, created_by
                FROM translated_files
                WHERE id IN \([$\d, ]+\);`

	e := errors.New("list err")
	mock.ExpectQuery(cmd).WillReturnRows().WillReturnError(e)

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

func TestPostgresqlTranslatedFileMetadataRepository_ListByIsid(t *testing.T) {
	mock, repo := newTranslMock(t)

	want := []*entity.TranslatedFileMetadata{{Id: 1, CreatedBy: "1"}, {Id: 2, CreatedBy: "1"}}

	columns := []string{"id", "original_files_id", "translated_filename", "target_language", "cost", "time_taken", "created_at", "updated_at", "created_by"}
	rows := sqlmock.NewRows(columns)
	for _, r := range want {
		rows.AddRow(r.Id, r.OriginalFileId, r.Filename, r.TargetLanguage, r.Cost, r.TimeTaken, r.CreatedAt, r.UpdatedAt, r.CreatedBy)
	}

	cmd := `SELECT id, original_files_id, translated_filename, target_language, cost, time_taken, created_by, updated_at, created_by
                FROM translated_files
                WHERE created_by = \$1;`

	mock.ExpectQuery(cmd).WillReturnRows(rows)

	got, err := repo.ListByIsid("1")
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

func TestPostgresqlTranslatedFileMetadataRepository_ListByIsid_Err(t *testing.T) {
	mock, repo := newTranslMock(t)

	cmd := `SELECT id, original_files_id, translated_filename, target_language, cost, time_taken, created_by, updated_at, created_by
                FROM translated_files
                WHERE created_by = \$1;`

	e := errors.New("list err")
	mock.ExpectQuery(cmd).WillReturnRows().WillReturnError(e)

	got, err := repo.ListByIsid("1")
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

func TestPostgresqlTranslatedFileMetadataRepository_ListOriginalFileIdsByIds(t *testing.T) {
	mock, repo := newTranslMock(t)

	want := []int{1, 2, 3}

	rows := sqlmock.NewRows([]string{"original_files_id"})
	for _, r := range want {
		rows.AddRow(r)
	}

	cmd := `SELECT original_files_id FROM translated_files WHERE id IN \([$\d, ]+\);`

	mock.ExpectQuery(cmd).WillReturnRows(rows)

	got, err := repo.ListOriginalFileIdsByIds([]int{1, 2})
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

func TestPostgresqlTranslatedFileMetadataRepository_ListOriginalFileIdsByIds_Err(t *testing.T) {
	mock, repo := newTranslMock(t)

	cmd := `SELECT original_files_id FROM translated_files WHERE id IN \([$\d, ]+\);`

	e := errors.New("list err")
	mock.ExpectQuery(cmd).WillReturnRows().WillReturnError(e)

	got, err := repo.ListOriginalFileIdsByIds([]int{1, 2})
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

func TestPostgresqlTranslatedFileMetadataRepository_Update(t *testing.T) {
	mock, repo := newTranslMock(t)

	ent := &entity.TranslatedFileMetadata{}

	cmd := `UPDATE translated_files SET
                original_files_id = \$1,
                translated_filename = \$2,
                target_language = \$3,
                cost = \$4,
                time_taken = \$5,
                created_at = \$6,
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

func TestPostgresqlTranslatedFileMetadataRepository_Update_Err(t *testing.T) {
	mock, repo := newTranslMock(t)

	ent := &entity.TranslatedFileMetadata{}

	cmd := `UPDATE translated_files SET
                original_files_id = \$1,
                translated_filename = \$2,
                target_language = \$3,
                cost = \$4,
                time_taken = \$5,
                created_at = \$6,
                updated_at = \$7,
                created_by = \$8
        WHERE id = \$9;`

	e := errors.New("update err")
	mock.ExpectExec(cmd).WillReturnResult(nil).WillReturnError(e)

	err := repo.Update(ent)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlTranslatedFileMetadataRepository_DeleteById(t *testing.T) {
	mock, repo := newTranslMock(t)

	cmd := `DELETE FROM translated_files WHERE id = \$1;`

	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(cmd).WillReturnResult(result)

	err := repo.DeleteById(1)
	if err != nil {
		t.Fatal(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlTranslatedFileMetadataRepository_DeleteById_Err(t *testing.T) {
	mock, repo := newTranslMock(t)

	cmd := `DELETE FROM translated_files WHERE id = \$1;`

	e := errors.New("delete err")
	mock.ExpectExec(cmd).WillReturnResult(nil).WillReturnError(e)

	err := repo.DeleteById(1)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlTranslatedFileMetadataRepository_DeleteByIds(t *testing.T) {
	mock, repo := newTranslMock(t)

	cmd := `DELETE FROM translated_files WHERE id IN \([$\d, ]+\);`

	result := sqlmock.NewResult(1, 2)
	mock.ExpectExec(cmd).WillReturnResult(result)

	err := repo.DeleteByIds([]int{1, 2})
	if err != nil {
		t.Fatal(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}

func TestPostgresqlTranslatedFileMetadataRepository_DeleteByIds_Err(t *testing.T) {
	mock, repo := newTranslMock(t)

	cmd := `DELETE FROM translated_files WHERE id IN \([$\d, ]+\);`

	e := errors.New("delete err")
	mock.ExpectExec(cmd).WillReturnResult(nil).WillReturnError(e)

	err := repo.DeleteByIds([]int{1, 2})
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(t)
	}
}
