package postgresql

import (
	"doc-translate-go/pkg/db"
	"doc-translate-go/pkg/user/entity"
	"doc-translate-go/pkg/user/repository"
)

type PostgresqlUserRepository struct {
	querier db.Querier
}

func NewPostgresqlUserRepository(querier db.Querier) *PostgresqlUserRepository {
	return &PostgresqlUserRepository{querier}
}

func (r *PostgresqlUserRepository) Create(u *entity.User) (int, error) {
	cmd := `INSERT INTO users (isid, role, email, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5)
                RETURNING id;`

	row := r.querier.QueryRow(cmd, u.Isid, u.Role, u.Email, u.CreatedAt, u.UpdatedAt)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresqlUserRepository) GetByIsid(isid string) (*entity.User, error) { panic("") }

func (r *PostgresqlUserRepository) Update(u *entity.User) error { panic("") }

func (r *PostgresqlUserRepository) DeleteById(id int) error { panic("") }

// Ensure implementation
var _ repository.UserRepository = (*PostgresqlUserRepository)(nil)
