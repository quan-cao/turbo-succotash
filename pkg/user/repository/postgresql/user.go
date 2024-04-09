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

func (r *PostgresqlUserRepository) GetByIsid(isid string) (*entity.User, error) {
	cmd := `SELECT id, isid, role, email, created_at, updated_at FROM users WHERE isid = $1;`

	row := r.querier.QueryRow(cmd, isid)

	var u entity.User
	err := row.Scan(&u.Id, &u.Isid, &u.Role, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresqlUserRepository) Update(u *entity.User) error {
	cmd := `UPDATE users SET isid = $1, role = $2, email = $3, updated_at = $4 WHERE id = $5;`
	_, err := r.querier.Exec(cmd, u.Isid, u.Role, u.Email, u.UpdatedAt, u.Id)
	return err
}

func (r *PostgresqlUserRepository) DeleteById(id int) error {
	cmd := `DELETE FROM users WHERE id = $1;`
	_, err := r.querier.Exec(cmd, id)
	return err
}

// Ensure implementation
var _ repository.UserRepository = (*PostgresqlUserRepository)(nil)
