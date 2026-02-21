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

// ==================== AcademicUnit ====================

type postgresAcademicUnitRepository struct{ db *sql.DB }

func NewPostgresAcademicUnitRepository(db *sql.DB) repository.AcademicUnitRepository {
	return &postgresAcademicUnitRepository{db: db}
}

const auCols = `id, parent_unit_id, school_id, name, code, type, description, level, academic_year, metadata, is_active, created_at, updated_at, deleted_at`

func scanAU(row interface{ Scan(...interface{}) error }) (*entities.AcademicUnit, error) {
	u := &entities.AcademicUnit{}
	err := row.Scan(&u.ID, &u.ParentUnitID, &u.SchoolID, &u.Name, &u.Code, &u.Type, &u.Description, &u.Level, &u.AcademicYear, &u.Metadata, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	return u, err
}

func (r *postgresAcademicUnitRepository) Create(ctx context.Context, unit *entities.AcademicUnit) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO academic_units (id, parent_unit_id, school_id, name, code, type, description, level, academic_year, metadata, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		unit.ID, unit.ParentUnitID, unit.SchoolID, unit.Name, unit.Code, unit.Type, unit.Description, unit.Level, unit.AcademicYear, unit.Metadata, unit.IsActive, unit.CreatedAt, unit.UpdatedAt)
	return err
}

func (r *postgresAcademicUnitRepository) FindByID(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error) {
	query := `SELECT ` + auCols + ` FROM academic_units WHERE id=$1`
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}
	u, err := scanAU(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *postgresAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	query := `SELECT ` + auCols + ` FROM academic_units WHERE school_id=$1`
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}
	query += ` ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var units []*entities.AcademicUnit
	for rows.Next() {
		u, err := scanAU(rows)
		if err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

func (r *postgresAcademicUnitRepository) FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	query := `SELECT ` + auCols + ` FROM academic_units WHERE school_id=$1 AND type=$2`
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}
	rows, err := r.db.QueryContext(ctx, query, schoolID, unitType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var units []*entities.AcademicUnit
	for rows.Next() {
		u, err := scanAU(rows)
		if err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

func (r *postgresAcademicUnitRepository) Update(ctx context.Context, unit *entities.AcademicUnit) error {
	_, err := r.db.ExecContext(ctx, `UPDATE academic_units SET parent_unit_id=$1, name=$2, description=$3, metadata=$4, updated_at=$5 WHERE id=$6`,
		unit.ParentUnitID, unit.Name, unit.Description, unit.Metadata, unit.UpdatedAt, unit.ID)
	return err
}

func (r *postgresAcademicUnitRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE academic_units SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND deleted_at IS NULL`, now, now, id)
	return err
}

func (r *postgresAcademicUnitRepository) Restore(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE academic_units SET deleted_at=NULL, is_active=true, updated_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}

func (r *postgresAcademicUnitRepository) GetHierarchyPath(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error) {
	query := `WITH RECURSIVE hierarchy AS (
		SELECT ` + auCols + ` FROM academic_units WHERE id=$1 AND deleted_at IS NULL
		UNION ALL
		SELECT au.id, au.parent_unit_id, au.school_id, au.name, au.code, au.type, au.description, au.level, au.academic_year, au.metadata, au.is_active, au.created_at, au.updated_at, au.deleted_at
		FROM academic_units au INNER JOIN hierarchy h ON au.id = h.parent_unit_id WHERE au.deleted_at IS NULL
	) SELECT * FROM hierarchy`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var units []*entities.AcademicUnit
	for rows.Next() {
		u, err := scanAU(rows)
		if err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	// Reverse to get root-first order
	for i, j := 0, len(units)-1; i < j; i, j = i+1, j-1 {
		units[i], units[j] = units[j], units[i]
	}
	return units, rows.Err()
}

func (r *postgresAcademicUnitRepository) ExistsBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM academic_units WHERE school_id=$1 AND code=$2 AND deleted_at IS NULL)`, schoolID, code).Scan(&exists)
	return exists, err
}

// ==================== Membership ====================

type postgresMembershipRepository struct{ db *sql.DB }

func NewPostgresMembershipRepository(db *sql.DB) repository.MembershipRepository {
	return &postgresMembershipRepository{db: db}
}

const mCols = `id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at`

func scanMembership(row interface{ Scan(...interface{}) error }) (*entities.Membership, error) {
	m := &entities.Membership{}
	err := row.Scan(&m.ID, &m.UserID, &m.SchoolID, &m.AcademicUnitID, &m.Role, &m.Metadata, &m.IsActive, &m.EnrolledAt, &m.WithdrawnAt, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (r *postgresMembershipRepository) Create(ctx context.Context, m *entities.Membership) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, m.ID, m.UserID, m.SchoolID, m.AcademicUnitID, m.Role, m.Metadata, m.IsActive, m.EnrolledAt, m.CreatedAt, m.UpdatedAt)
	return err
}

func (r *postgresMembershipRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	m, err := scanMembership(r.db.QueryRowContext(ctx, `SELECT `+mCols+` FROM memberships WHERE id=$1`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *postgresMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error) {
	return r.queryMemberships(ctx, `SELECT `+mCols+` FROM memberships WHERE user_id=$1 AND is_active=true ORDER BY created_at DESC`, userID)
}

func (r *postgresMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error) {
	return r.queryMemberships(ctx, `SELECT `+mCols+` FROM memberships WHERE academic_unit_id=$1 AND is_active=true ORDER BY created_at DESC`, unitID)
}

func (r *postgresMembershipRepository) FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error) {
	query := `SELECT ` + mCols + ` FROM memberships WHERE academic_unit_id=$1 AND role=$2`
	if activeOnly {
		query += ` AND is_active=true`
	}
	return r.queryMemberships(ctx, query, unitID, role)
}

func (r *postgresMembershipRepository) FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error) {
	m, err := scanMembership(r.db.QueryRowContext(ctx, `SELECT `+mCols+` FROM memberships WHERE user_id=$1 AND school_id=$2 AND is_active=true LIMIT 1`, userID, schoolID))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *postgresMembershipRepository) Update(ctx context.Context, m *entities.Membership) error {
	_, err := r.db.ExecContext(ctx, `UPDATE memberships SET role=$1, is_active=$2, withdrawn_at=$3, updated_at=$4 WHERE id=$5`,
		m.Role, m.IsActive, m.WithdrawnAt, m.UpdatedAt, m.ID)
	return err
}

func (r *postgresMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM memberships WHERE id=$1`, id)
	return err
}

func (r *postgresMembershipRepository) queryMemberships(ctx context.Context, query string, args ...interface{}) ([]*entities.Membership, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*entities.Membership
	for rows.Next() {
		m, err := scanMembership(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

// ==================== Role ====================

type postgresRoleRepository struct{ db *sql.DB }

func NewPostgresRoleRepository(db *sql.DB) repository.RoleRepository {
	return &postgresRoleRepository{db: db}
}

func scanRole(row interface{ Scan(...interface{}) error }) (*entities.Role, error) {
	r := &entities.Role{}
	err := row.Scan(&r.ID, &r.Name, &r.DisplayName, &r.Description, &r.Scope, &r.IsActive, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func (r *postgresRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	role, err := scanRole(r.db.QueryRowContext(ctx, `SELECT id, name, display_name, description, scope, is_active, created_at, updated_at FROM roles WHERE id=$1 AND is_active=true`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return role, err
}

func (r *postgresRoleRepository) FindAll(ctx context.Context) ([]*entities.Role, error) {
	return r.queryRoles(ctx, `SELECT id, name, display_name, description, scope, is_active, created_at, updated_at FROM roles WHERE is_active=true ORDER BY name`)
}

func (r *postgresRoleRepository) FindByScope(ctx context.Context, scope string) ([]*entities.Role, error) {
	return r.queryRoles(ctx, `SELECT id, name, display_name, description, scope, is_active, created_at, updated_at FROM roles WHERE scope=$1 AND is_active=true ORDER BY name`, scope)
}

func (r *postgresRoleRepository) queryRoles(ctx context.Context, query string, args ...interface{}) ([]*entities.Role, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []*entities.Role
	for rows.Next() {
		role, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

// ==================== Permission ====================

type postgresPermissionRepository struct{ db *sql.DB }

func NewPostgresPermissionRepository(db *sql.DB) repository.PermissionRepository {
	return &postgresPermissionRepository{db: db}
}

const permCols = `p.id, p.name, p.display_name, p.description, p.resource_id, COALESCE(r.key, ''), p.action, p.scope, p.is_active, p.created_at, p.updated_at`

func scanPerm(row interface{ Scan(...interface{}) error }) (*entities.Permission, error) {
	p := &entities.Permission{}
	err := row.Scan(&p.ID, &p.Name, &p.DisplayName, &p.Description, &p.ResourceID, &p.ResourceKey, &p.Action, &p.Scope, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *postgresPermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	p, err := scanPerm(r.db.QueryRowContext(ctx, `SELECT `+permCols+` FROM permissions p LEFT JOIN resources r ON p.resource_id=r.id WHERE p.id=$1 AND p.is_active=true`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *postgresPermissionRepository) FindAll(ctx context.Context) ([]*entities.Permission, error) {
	return r.queryPerms(ctx, `SELECT `+permCols+` FROM permissions p LEFT JOIN resources r ON p.resource_id=r.id WHERE p.is_active=true ORDER BY p.name`)
}

func (r *postgresPermissionRepository) FindByRole(ctx context.Context, roleID uuid.UUID) ([]*entities.Permission, error) {
	return r.queryPerms(ctx, `SELECT `+permCols+` FROM permissions p LEFT JOIN resources r ON p.resource_id=r.id
		INNER JOIN role_permissions rp ON p.id=rp.permission_id WHERE rp.role_id=$1 AND p.is_active=true ORDER BY p.name`, roleID)
}

func (r *postgresPermissionRepository) queryPerms(ctx context.Context, query string, args ...interface{}) ([]*entities.Permission, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []*entities.Permission
	for rows.Next() {
		p, err := scanPerm(rows)
		if err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

// ==================== UserRole ====================

type postgresUserRoleRepository struct{ db *sql.DB }

func NewPostgresUserRoleRepository(db *sql.DB) repository.UserRoleRepository {
	return &postgresUserRoleRepository{db: db}
}

func scanUserRole(row interface{ Scan(...interface{}) error }) (*entities.UserRole, error) {
	ur := &entities.UserRole{}
	err := row.Scan(&ur.ID, &ur.UserID, &ur.RoleID, &ur.SchoolID, &ur.AcademicUnitID, &ur.IsActive, &ur.GrantedBy, &ur.GrantedAt, &ur.ExpiresAt, &ur.CreatedAt, &ur.UpdatedAt)
	return ur, err
}

const urCols = `id, user_id, role_id, school_id, academic_unit_id, is_active, granted_by, granted_at, expires_at, created_at, updated_at`

func (r *postgresUserRoleRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.UserRole, error) {
	return r.queryUserRoles(ctx, `SELECT `+urCols+` FROM user_roles WHERE user_id=$1 AND is_active=true`, userID)
}

func (r *postgresUserRoleRepository) FindByUserInContext(ctx context.Context, userID uuid.UUID, schoolID *uuid.UUID, unitID *uuid.UUID) ([]*entities.UserRole, error) {
	query := `SELECT ` + urCols + ` FROM user_roles WHERE user_id=$1 AND is_active=true`
	args := []interface{}{userID}
	argN := 2
	if schoolID != nil {
		query += fmt.Sprintf(` AND school_id=$%d`, argN)
		args = append(args, *schoolID)
		argN++
	}
	if unitID != nil {
		query += fmt.Sprintf(` AND academic_unit_id=$%d`, argN)
		args = append(args, *unitID)
	}
	query += ` ORDER BY created_at`
	return r.queryUserRoles(ctx, query, args...)
}

func (r *postgresUserRoleRepository) Grant(ctx context.Context, ur *entities.UserRole) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO user_roles (id, user_id, role_id, school_id, academic_unit_id, is_active, granted_by, granted_at, expires_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		ur.ID, ur.UserID, ur.RoleID, ur.SchoolID, ur.AcademicUnitID, ur.IsActive, ur.GrantedBy, ur.GrantedAt, ur.ExpiresAt, ur.CreatedAt, ur.UpdatedAt)
	return err
}

func (r *postgresUserRoleRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE user_roles SET is_active=false, updated_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}

func (r *postgresUserRoleRepository) RevokeByUserAndRole(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) error {
	query := `UPDATE user_roles SET is_active=false, updated_at=$1 WHERE user_id=$2 AND role_id=$3 AND is_active=true`
	args := []interface{}{time.Now(), userID, roleID}
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *postgresUserRoleRepository) UserHasRole(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id=$1 AND role_id=$2 AND is_active=true`
	args := []interface{}{userID, roleID}
	argN := 3
	if schoolID != nil {
		query += fmt.Sprintf(` AND school_id=$%d`, argN)
		args = append(args, *schoolID)
		argN++
	}
	if unitID != nil {
		query += fmt.Sprintf(` AND academic_unit_id=$%d`, argN)
		args = append(args, *unitID)
	}
	query += `)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	return exists, err
}

func (r *postgresUserRoleRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID, schoolID, unitID *uuid.UUID) ([]string, error) {
	query := `SELECT DISTINCT p.name FROM permissions p
		INNER JOIN role_permissions rp ON p.id=rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id=ur.role_id
		WHERE ur.user_id=$1 AND ur.is_active=true AND p.is_active=true`
	args := []interface{}{userID}
	argN := 2
	if schoolID != nil {
		query += fmt.Sprintf(` AND ur.school_id=$%d`, argN)
		args = append(args, *schoolID)
		argN++
	}
	if unitID != nil {
		query += fmt.Sprintf(` AND ur.academic_unit_id=$%d`, argN)
		args = append(args, *unitID)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, rows.Err()
}

func (r *postgresUserRoleRepository) queryUserRoles(ctx context.Context, query string, args ...interface{}) ([]*entities.UserRole, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*entities.UserRole
	for rows.Next() {
		ur, err := scanUserRole(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, ur)
	}
	return result, rows.Err()
}

// ==================== Resource ====================

type postgresResourceRepository struct{ db *sql.DB }

func NewPostgresResourceRepository(db *sql.DB) repository.ResourceRepository {
	return &postgresResourceRepository{db: db}
}

const resCols = `id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope, is_active, created_at, updated_at`

func scanResource(row interface{ Scan(...interface{}) error }) (*entities.Resource, error) {
	r := &entities.Resource{}
	err := row.Scan(&r.ID, &r.Key, &r.DisplayName, &r.Description, &r.Icon, &r.ParentID, &r.SortOrder, &r.IsMenuVisible, &r.Scope, &r.IsActive, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func (r *postgresResourceRepository) FindAll(ctx context.Context) ([]*entities.Resource, error) {
	return r.queryResources(ctx, `SELECT `+resCols+` FROM resources WHERE is_active=true ORDER BY sort_order`)
}

func (r *postgresResourceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Resource, error) {
	res, err := scanResource(r.db.QueryRowContext(ctx, `SELECT `+resCols+` FROM resources WHERE id=$1 AND is_active=true`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return res, err
}

func (r *postgresResourceRepository) FindMenuVisible(ctx context.Context) ([]*entities.Resource, error) {
	return r.queryResources(ctx, `SELECT `+resCols+` FROM resources WHERE is_menu_visible=true AND is_active=true ORDER BY sort_order`)
}

func (r *postgresResourceRepository) Create(ctx context.Context, res *entities.Resource) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		res.ID, res.Key, res.DisplayName, res.Description, res.Icon, res.ParentID, res.SortOrder, res.IsMenuVisible, res.Scope, res.IsActive, res.CreatedAt, res.UpdatedAt)
	return err
}

func (r *postgresResourceRepository) Update(ctx context.Context, res *entities.Resource) error {
	_, err := r.db.ExecContext(ctx, `UPDATE resources SET display_name=$1, description=$2, icon=$3, parent_id=$4, sort_order=$5, is_menu_visible=$6, scope=$7, is_active=$8, updated_at=$9 WHERE id=$10`,
		res.DisplayName, res.Description, res.Icon, res.ParentID, res.SortOrder, res.IsMenuVisible, res.Scope, res.IsActive, res.UpdatedAt, res.ID)
	return err
}

func (r *postgresResourceRepository) queryResources(ctx context.Context, query string, args ...interface{}) ([]*entities.Resource, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*entities.Resource
	for rows.Next() {
		res, err := scanResource(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, rows.Err()
}

// ==================== Subject ====================

type postgresSubjectRepository struct{ db *sql.DB }

func NewPostgresSubjectRepository(db *sql.DB) repository.SubjectRepository {
	return &postgresSubjectRepository{db: db}
}

func scanSubject(row interface{ Scan(...interface{}) error }) (*entities.Subject, error) {
	s := &entities.Subject{}
	err := row.Scan(&s.ID, &s.Name, &s.Description, &s.Metadata, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *postgresSubjectRepository) Create(ctx context.Context, s *entities.Subject) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO subjects (id, name, description, metadata, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		s.ID, s.Name, s.Description, s.Metadata, s.IsActive, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *postgresSubjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Subject, error) {
	s, err := scanSubject(r.db.QueryRowContext(ctx, `SELECT id, name, description, metadata, is_active, created_at, updated_at FROM subjects WHERE id=$1 AND is_active=true`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *postgresSubjectRepository) Update(ctx context.Context, s *entities.Subject) error {
	_, err := r.db.ExecContext(ctx, `UPDATE subjects SET name=$1, description=$2, metadata=$3, updated_at=$4 WHERE id=$5`, s.Name, s.Description, s.Metadata, s.UpdatedAt, s.ID)
	return err
}

func (r *postgresSubjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE subjects SET is_active=false, updated_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}

func (r *postgresSubjectRepository) List(ctx context.Context) ([]*entities.Subject, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, metadata, is_active, created_at, updated_at FROM subjects WHERE is_active=true ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subjects []*entities.Subject
	for rows.Next() {
		s, err := scanSubject(rows)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, rows.Err()
}

func (r *postgresSubjectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM subjects WHERE name=$1 AND is_active=true)`, name).Scan(&exists)
	return exists, err
}

// ==================== Guardian ====================

type postgresGuardianRepository struct{ db *sql.DB }

func NewPostgresGuardianRepository(db *sql.DB) repository.GuardianRepository {
	return &postgresGuardianRepository{db: db}
}

const grCols = `id, guardian_id, student_id, relationship_type, is_active, created_at, updated_at, created_by`

func scanGR(row interface{ Scan(...interface{}) error }) (*entities.GuardianRelation, error) {
	g := &entities.GuardianRelation{}
	err := row.Scan(&g.ID, &g.GuardianID, &g.StudentID, &g.RelationshipType, &g.IsActive, &g.CreatedAt, &g.UpdatedAt, &g.CreatedBy)
	return g, err
}

func (r *postgresGuardianRepository) Create(ctx context.Context, g *entities.GuardianRelation) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO guardian_relations (id, guardian_id, student_id, relationship_type, is_active, created_at, updated_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		g.ID, g.GuardianID, g.StudentID, g.RelationshipType, g.IsActive, g.CreatedAt, g.UpdatedAt, g.CreatedBy)
	return err
}

func (r *postgresGuardianRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	g, err := scanGR(r.db.QueryRowContext(ctx, `SELECT `+grCols+` FROM guardian_relations WHERE id=$1`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return g, err
}

func (r *postgresGuardianRepository) FindByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	return r.queryGR(ctx, `SELECT `+grCols+` FROM guardian_relations WHERE guardian_id=$1 AND is_active=true`, guardianID)
}

func (r *postgresGuardianRepository) FindByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	return r.queryGR(ctx, `SELECT `+grCols+` FROM guardian_relations WHERE student_id=$1 AND is_active=true`, studentID)
}

func (r *postgresGuardianRepository) Update(ctx context.Context, g *entities.GuardianRelation) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guardian_relations SET relationship_type=$1, is_active=$2, updated_at=$3 WHERE id=$4`,
		g.RelationshipType, g.IsActive, g.UpdatedAt, g.ID)
	return err
}

func (r *postgresGuardianRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guardian_relations SET is_active=false, updated_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}

func (r *postgresGuardianRepository) ExistsActiveRelation(ctx context.Context, guardianID, studentID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM guardian_relations WHERE guardian_id=$1 AND student_id=$2 AND is_active=true)`, guardianID, studentID).Scan(&exists)
	return exists, err
}

func (r *postgresGuardianRepository) queryGR(ctx context.Context, query string, args ...interface{}) ([]*entities.GuardianRelation, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*entities.GuardianRelation
	for rows.Next() {
		g, err := scanGR(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	return result, rows.Err()
}

// ==================== ScreenTemplate ====================

type postgresScreenTemplateRepository struct{ db *sql.DB }

func NewPostgresScreenTemplateRepository(db *sql.DB) repository.ScreenTemplateRepository {
	return &postgresScreenTemplateRepository{db: db}
}

func (r *postgresScreenTemplateRepository) Create(ctx context.Context, t *entities.ScreenTemplate) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO ui_config.screen_templates (id, pattern, name, description, version, definition, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		t.ID, t.Pattern, t.Name, t.Description, t.Version, t.Definition, t.IsActive, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *postgresScreenTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ScreenTemplate, error) {
	t := &entities.ScreenTemplate{}
	err := r.db.QueryRowContext(ctx, `SELECT id, pattern, name, description, version, definition, is_active, created_at, updated_at FROM ui_config.screen_templates WHERE id=$1`, id).
		Scan(&t.ID, &t.Pattern, &t.Name, &t.Description, &t.Version, &t.Definition, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("screen template not found")
	}
	return t, err
}

func (r *postgresScreenTemplateRepository) List(ctx context.Context, filter repository.ScreenTemplateFilter) ([]*entities.ScreenTemplate, int, error) {
	countQuery := `SELECT COUNT(*) FROM ui_config.screen_templates WHERE is_active=true`
	query := `SELECT id, pattern, name, description, version, definition, is_active, created_at, updated_at FROM ui_config.screen_templates WHERE is_active=true`
	args := []interface{}{}
	argN := 1
	if filter.Pattern != "" {
		countQuery += fmt.Sprintf(` AND pattern=$%d`, argN)
		query += fmt.Sprintf(` AND pattern=$%d`, argN)
		args = append(args, filter.Pattern)
		argN++
	}
	var total int
	_ = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)

	query += ` ORDER BY created_at DESC`
	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, argN)
		args = append(args, filter.Limit)
		argN++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET $%d`, argN)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var templates []*entities.ScreenTemplate
	for rows.Next() {
		t := &entities.ScreenTemplate{}
		if err := rows.Scan(&t.ID, &t.Pattern, &t.Name, &t.Description, &t.Version, &t.Definition, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		templates = append(templates, t)
	}
	return templates, total, rows.Err()
}

func (r *postgresScreenTemplateRepository) Update(ctx context.Context, t *entities.ScreenTemplate) error {
	_, err := r.db.ExecContext(ctx, `UPDATE ui_config.screen_templates SET pattern=$1, name=$2, description=$3, version=$4, definition=$5, updated_at=$6 WHERE id=$7`,
		t.Pattern, t.Name, t.Description, t.Version, t.Definition, t.UpdatedAt, t.ID)
	return err
}

func (r *postgresScreenTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE ui_config.screen_templates SET is_active=false, updated_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}

// ==================== ScreenInstance ====================

type postgresScreenInstanceRepository struct{ db *sql.DB }

func NewPostgresScreenInstanceRepository(db *sql.DB) repository.ScreenInstanceRepository {
	return &postgresScreenInstanceRepository{db: db}
}

func (r *postgresScreenInstanceRepository) Create(ctx context.Context, i *entities.ScreenInstance) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		i.ID, i.ScreenKey, i.TemplateID, i.Name, i.Description, i.SlotData, i.Actions, i.DataEndpoint, i.DataConfig, i.Scope, i.RequiredPermission, i.HandlerKey, i.IsActive, i.CreatedAt, i.UpdatedAt)
	return err
}

func (r *postgresScreenInstanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ScreenInstance, error) {
	i := &entities.ScreenInstance{}
	err := r.db.QueryRowContext(ctx, `SELECT id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key, is_active, created_at, updated_at
		FROM ui_config.screen_instances WHERE id=$1`, id).Scan(&i.ID, &i.ScreenKey, &i.TemplateID, &i.Name, &i.Description, &i.SlotData, &i.Actions, &i.DataEndpoint, &i.DataConfig, &i.Scope, &i.RequiredPermission, &i.HandlerKey, &i.IsActive, &i.CreatedAt, &i.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("screen instance not found")
	}
	return i, err
}

func (r *postgresScreenInstanceRepository) GetByScreenKey(ctx context.Context, key string) (*entities.ScreenInstance, error) {
	i := &entities.ScreenInstance{}
	err := r.db.QueryRowContext(ctx, `SELECT id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key, is_active, created_at, updated_at
		FROM ui_config.screen_instances WHERE screen_key=$1 AND is_active=true`, key).Scan(&i.ID, &i.ScreenKey, &i.TemplateID, &i.Name, &i.Description, &i.SlotData, &i.Actions, &i.DataEndpoint, &i.DataConfig, &i.Scope, &i.RequiredPermission, &i.HandlerKey, &i.IsActive, &i.CreatedAt, &i.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("screen instance not found for key: %s", key)
	}
	return i, err
}

func (r *postgresScreenInstanceRepository) List(ctx context.Context, filter repository.ScreenInstanceFilter) ([]*entities.ScreenInstance, int, error) {
	countQuery := `SELECT COUNT(*) FROM ui_config.screen_instances WHERE is_active=true`
	query := `SELECT id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key, is_active, created_at, updated_at FROM ui_config.screen_instances WHERE is_active=true`
	args := []interface{}{}
	argN := 1
	if filter.TemplateID != nil {
		countQuery += fmt.Sprintf(` AND template_id=$%d`, argN)
		query += fmt.Sprintf(` AND template_id=$%d`, argN)
		args = append(args, *filter.TemplateID)
		argN++
	}
	var total int
	_ = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	query += ` ORDER BY created_at DESC`
	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, argN)
		args = append(args, filter.Limit)
		argN++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET $%d`, argN)
		args = append(args, filter.Offset)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var instances []*entities.ScreenInstance
	for rows.Next() {
		i := &entities.ScreenInstance{}
		if err := rows.Scan(&i.ID, &i.ScreenKey, &i.TemplateID, &i.Name, &i.Description, &i.SlotData, &i.Actions, &i.DataEndpoint, &i.DataConfig, &i.Scope, &i.RequiredPermission, &i.HandlerKey, &i.IsActive, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, 0, err
		}
		instances = append(instances, i)
	}
	return instances, total, rows.Err()
}

func (r *postgresScreenInstanceRepository) Update(ctx context.Context, i *entities.ScreenInstance) error {
	_, err := r.db.ExecContext(ctx, `UPDATE ui_config.screen_instances SET screen_key=$1, template_id=$2, name=$3, description=$4, slot_data=$5, actions=$6, data_endpoint=$7, data_config=$8, scope=$9, required_permission=$10, handler_key=$11, updated_at=$12 WHERE id=$13`,
		i.ScreenKey, i.TemplateID, i.Name, i.Description, i.SlotData, i.Actions, i.DataEndpoint, i.DataConfig, i.Scope, i.RequiredPermission, i.HandlerKey, i.UpdatedAt, i.ID)
	return err
}

func (r *postgresScreenInstanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE ui_config.screen_instances SET is_active=false, updated_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}

// ==================== ResourceScreen ====================

type postgresResourceScreenRepository struct{ db *sql.DB }

func NewPostgresResourceScreenRepository(db *sql.DB) repository.ResourceScreenRepository {
	return &postgresResourceScreenRepository{db: db}
}

func (r *postgresResourceScreenRepository) Create(ctx context.Context, rs *entities.ResourceScreen) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, rs.ID, rs.ResourceID, rs.ResourceKey, rs.ScreenKey, rs.ScreenType, rs.IsDefault, rs.SortOrder, rs.IsActive, rs.CreatedAt, rs.UpdatedAt)
	return err
}

func (r *postgresResourceScreenRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*entities.ResourceScreen, error) {
	return r.queryRS(ctx, `SELECT id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order, is_active, created_at, updated_at FROM ui_config.resource_screens WHERE resource_id=$1 AND is_active=true ORDER BY sort_order`, resourceID)
}

func (r *postgresResourceScreenRepository) GetByResourceKey(ctx context.Context, key string) ([]*entities.ResourceScreen, error) {
	return r.queryRS(ctx, `SELECT id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order, is_active, created_at, updated_at FROM ui_config.resource_screens WHERE resource_key=$1 AND is_active=true ORDER BY sort_order`, key)
}

func (r *postgresResourceScreenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ui_config.resource_screens WHERE id=$1`, id)
	return err
}

func (r *postgresResourceScreenRepository) queryRS(ctx context.Context, query string, args ...interface{}) ([]*entities.ResourceScreen, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*entities.ResourceScreen
	for rows.Next() {
		rs := &entities.ResourceScreen{}
		if err := rows.Scan(&rs.ID, &rs.ResourceID, &rs.ResourceKey, &rs.ScreenKey, &rs.ScreenType, &rs.IsDefault, &rs.SortOrder, &rs.IsActive, &rs.CreatedAt, &rs.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, rs)
	}
	return result, rows.Err()
}
