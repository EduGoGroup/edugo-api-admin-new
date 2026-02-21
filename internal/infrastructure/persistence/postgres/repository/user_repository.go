package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type postgresUserRepository struct{ db *sql.DB }

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

const userCols = `id, email, password_hash, first_name, last_name, school_id, is_active, created_at, updated_at, deleted_at`

func scanUser(row interface{ Scan(...interface{}) error }) (*entities.User, error) {
	u := &entities.User{}
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.SchoolID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	return u, err
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	u, err := scanUser(r.db.QueryRowContext(ctx, `SELECT `+userCols+` FROM users WHERE id=$1 AND deleted_at IS NULL`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	u, err := scanUser(r.db.QueryRowContext(ctx, `SELECT `+userCols+` FROM users WHERE email=$1 AND deleted_at IS NULL`, email))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *postgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET first_name=$1, last_name=$2, is_active=$3, updated_at=$4 WHERE id=$5 AND deleted_at IS NULL`,
		user.FirstName, user.LastName, user.IsActive, user.UpdatedAt, user.ID)
	return err
}

func (r *postgresUserRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.User, error) {
	query := `SELECT ` + userCols + ` FROM users WHERE deleted_at IS NULL`
	args := []interface{}{}
	argN := 1
	if filters.IsActive != nil {
		query += fmt.Sprintf(` AND is_active=$%d`, argN)
		args = append(args, *filters.IsActive)
		argN++
	}
	query += ` ORDER BY created_at DESC`
	if filters.Limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, argN)
		args = append(args, filters.Limit)
		argN++
	}
	if filters.Offset > 0 {
		query += fmt.Sprintf(` OFFSET $%d`, argN)
		args = append(args, filters.Offset)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*entities.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
