package postgresql

import (
	"database/sql"
	"doc-translate-go/pkg/user/entity"
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

	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)

	got, err := repo.Create(ent)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != 0 {
		t.Fatalf("expected %v, got %v", 0, got)
	}
}
