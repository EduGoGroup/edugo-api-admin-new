package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type postgresSchoolRepository struct{ db *sql.DB }

func NewPostgresSchoolRepository(db *sql.DB) repository.SchoolRepository {
	return &postgresSchoolRepository{db: db}
}

func (r *postgresSchoolRepository) Create(ctx context.Context, school *entities.School) error {
	query := `INSERT INTO schools (id, name, code, address, city, country, phone, email, metadata,
		is_active, subscription_tier, max_teachers, max_students, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.db.ExecContext(ctx, query, school.ID, school.Name, school.Code, school.Address, school.City,
		school.Country, school.Phone, school.Email, school.Metadata, school.IsActive,
		school.SubscriptionTier, school.MaxTeachers, school.MaxStudents, school.CreatedAt, school.UpdatedAt)
	return err
}

const schoolColumns = `id, name, code, address, city, country, phone, email, metadata,
	is_active, subscription_tier, max_teachers, max_students, created_at, updated_at, deleted_at`

func scanSchool(row interface{ Scan(...interface{}) error }) (*entities.School, error) {
	s := &entities.School{}
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.Address, &s.City, &s.Country, &s.Phone, &s.Email,
		&s.Metadata, &s.IsActive, &s.SubscriptionTier, &s.MaxTeachers, &s.MaxStudents, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
	return s, err
}

func (r *postgresSchoolRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error) {
	query := `SELECT ` + schoolColumns + ` FROM schools WHERE id = $1 AND deleted_at IS NULL`
	s, err := scanSchool(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *postgresSchoolRepository) FindByCode(ctx context.Context, code string) (*entities.School, error) {
	query := `SELECT ` + schoolColumns + ` FROM schools WHERE code = $1 AND deleted_at IS NULL`
	s, err := scanSchool(r.db.QueryRowContext(ctx, query, code))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *postgresSchoolRepository) Update(ctx context.Context, school *entities.School) error {
	query := `UPDATE schools SET name=$1, address=$2, city=$3, country=$4, phone=$5, email=$6,
		metadata=$7, is_active=$8, subscription_tier=$9, max_teachers=$10, max_students=$11, updated_at=$12
		WHERE id=$13 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, school.Name, school.Address, school.City, school.Country,
		school.Phone, school.Email, school.Metadata, school.IsActive, school.SubscriptionTier,
		school.MaxTeachers, school.MaxStudents, school.UpdatedAt, school.ID)
	return err
}

func (r *postgresSchoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE schools SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND deleted_at IS NULL`, now, now, id)
	return err
}

func (r *postgresSchoolRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.School, error) {
	query := `SELECT ` + schoolColumns + ` FROM schools WHERE deleted_at IS NULL`
	args := []interface{}{}
	argN := 1

	if filters.IsActive != nil {
		query += fmt.Sprintf(` AND is_active = $%d`, argN)
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

	var schools []*entities.School
	for rows.Next() {
		s, err := scanSchool(rows)
		if err != nil {
			return nil, err
		}
		schools = append(schools, s)
	}
	return schools, rows.Err()
}

func (r *postgresSchoolRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM schools WHERE code=$1 AND deleted_at IS NULL)`, code).Scan(&exists)
	return exists, err
}
