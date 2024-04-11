package postgresql

import (
	"database/sql"
	"doc-translate-go/pkg/user/entity"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresqlOriginalFileMetadataRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{Id: 1}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(ent.Id)

	cmd := `INSERT INTO users \(isid, role, email, created_at, updated_at\)
                VALUES \(\$1, \$2, \$3, \$4, \$5\)
                RETURNING id;`

	mock.ExpectQuery(cmd).WillReturnRows(rows)

	got, err := repo.Create(ent)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if got != ent.Id {
		t.Fatalf("expected %v, got %v", ent.Id, got)
	}

	mock.ExpectationsWereMet()
}

func TestPostgresqlOriginalFileMetadataRepository_Create_ScanErr(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{Id: 1}

	cmd := `INSERT INTO users \(isid, role, email, created_at, updated_at\)
                VALUES \(\$1, \$2, \$3, \$4, \$5\)
                RETURNING id;`

	e := sql.ErrNoRows
	rows := sqlmock.NewRows([]string{"id"})

	mock.ExpectQuery(cmd).WillReturnRows(rows).WillReturnError(e)

	got, err := repo.Create(ent)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != 0 {
		t.Fatalf("expected %v, got %v", 0, got)
	}

	mock.ExpectationsWereMet()
}

func TestPostgresqlOriginalFileMetadataRepository_GetByIsid(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{Isid: "1"}

	rows := sqlmock.NewRows([]string{"id", "isid", "role", "email", "created_at", "updated_at"}).
		AddRow(ent.Id, ent.Isid, ent.Role, ent.Email, ent.CreatedAt, ent.UpdatedAt)

	cmd := `SELECT id, isid, role, email, created_at, updated_at FROM users WHERE isid = \$1;`

	mock.ExpectQuery(cmd).WillReturnRows(rows)

	got, err := repo.GetByIsid(ent.Isid)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if !reflect.DeepEqual(got, ent) {
		t.Fatalf("expected %v, got %v", ent, got)
	}

	mock.ExpectationsWereMet()
}

func TestPostgresqlOriginalFileMetadataRepository_GetByIsid_ScanErr(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{Isid: "1"}

	rows := sqlmock.NewRows([]string{"id", "isid", "role", "email", "created_at", "updated_at"})

	cmd := `SELECT id, isid, role, email, created_at, updated_at FROM users WHERE isid = \$1;`

	e := sql.ErrNoRows
	mock.ExpectQuery(cmd).WillReturnRows(rows).WillReturnError(e)

	got, err := repo.GetByIsid(ent.Isid)
	if err != e {
		t.Fatalf("expected nil, got %v", err)
	}

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	mock.ExpectationsWereMet()
}

func TestPostgresqlOriginalFileMetadataRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{}

	cmd := `UPDATE users SET isid = \$1, role = \$2, email = \$3, updated_at = \$4 WHERE id = \$5;`

	mock.
		ExpectExec(cmd).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Update(ent); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	mock.ExpectationsWereMet()
}

func TestPostgresqlOriginalFileMetadataRepository_Update_Err(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{}

	cmd := `UPDATE users SET isid = \$1, role = \$2, email = \$3, updated_at = \$4 WHERE id = \$5;`

	e := errors.New("exec error")
	mock.
		ExpectExec(cmd).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(e)

	if err := repo.Update(ent); err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	mock.ExpectationsWereMet()
}

func TestPostgresqlOriginalFileMetadataRepository_DeleteById(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{}

	cmd := `DELETE FROM users WHERE id = \$1;`

	mock.
		ExpectExec(cmd).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.DeleteById(ent.Id); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	mock.ExpectationsWereMet()
}

func TestPostgresqlOriginalFileMetadataRepository_DeleteById_Err(t *testing.T) {
	db, mock, _ := sqlmock.New()

	repo := NewPostgresqlUserRepository(db)
	ent := &entity.User{}

	cmd := `DELETE FROM users WHERE id = \$1;`

	e := errors.New("exec error")
	mock.
		ExpectExec(cmd).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(e)

	if err := repo.DeleteById(ent.Id); err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	mock.ExpectationsWereMet()
}
