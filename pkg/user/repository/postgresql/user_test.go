package postgresql

import (
	"database/sql"
	"doc-translate-go/pkg/user/entity"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresqlOriginalFileMetadataRepository__Create(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{Id: 1}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(ent.Id)

	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)

	got, err := repo.Create(ent)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if got != ent.Id {
		t.Fatalf("expected %v, got %v", ent.Id, got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__Create_ScanErr(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{Id: 1}

	e := sql.ErrNoRows
	rows := sqlmock.NewRows([]string{"id"})

	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows).WillReturnError(e)

	got, err := repo.Create(ent)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != 0 {
		t.Fatalf("expected %v, got %v", 0, got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__GetByIsid(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{Isid: "1"}

	rows := sqlmock.NewRows([]string{"id", "isid", "role", "email", "created_at", "updated_at"}).
		AddRow(ent.Id, ent.Isid, ent.Role, ent.Email, ent.CreatedAt, ent.UpdatedAt)

	mock.ExpectQuery("SELECT .* FROM users WHERE isid = \\$1;").WillReturnRows(rows)

	got, err := repo.GetByIsid(ent.Isid)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if !reflect.DeepEqual(got, ent) {
		t.Fatalf("expected %v, got %v", ent, got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__GetByIsid_ScanErr(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{Isid: "1"}

	rows := sqlmock.NewRows([]string{"id", "isid", "role", "email", "created_at", "updated_at"})

	e := sql.ErrNoRows
	mock.ExpectQuery("SELECT .* FROM users WHERE isid = \\$1;").WillReturnRows(rows).WillReturnError(e)

	got, err := repo.GetByIsid(ent.Isid)
	if err != e {
		t.Fatalf("expected nil, got %v", err)
	}

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__Update(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{}

	mock.
		ExpectExec("UPDATE users SET isid = \\$1, role = \\$2, email = \\$3, updated_at = \\$4 WHERE id = \\$5;").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Update(ent); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__Update_Err(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{}

	e := errors.New("exec error")
	mock.
		ExpectExec("UPDATE users SET isid = \\$1, role = \\$2, email = \\$3, updated_at = \\$4 WHERE id = \\$5;").
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(e)

	if err := repo.Update(ent); err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__DeleteById(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{}

	mock.
		ExpectExec("DELETE FROM users WHERE id = \\$1;").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.DeleteById(ent.Id); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestPostgresqlOriginalFileMetadataRepository__DeleteById_Err(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{}

	e := errors.New("exec error")
	mock.
		ExpectExec("DELETE FROM users WHERE id = \\$1;").
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(e)

	if err := repo.DeleteById(ent.Id); err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}
