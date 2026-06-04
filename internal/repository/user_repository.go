package repository

import (
	"context"
	"errors"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (id, employee_code, full_name, email, phone, password_hash, role, department, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`
	
	return r.db.QueryRow(context.Background(), query,
		user.ID, user.EmployeeCode, user.FullName, user.Email, user.Phone,
		user.PasswordHash, user.Role, user.Department, user.Status,
		user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
}

func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	query := `SELECT id, employee_code, full_name, email, phone, password_hash, role, department, status, created_at, updated_at FROM users WHERE phone = $1`
	
	err := r.db.QueryRow(context.Background(), query, phone).Scan(
		&user.ID, &user.EmployeeCode, &user.FullName, &user.Email, &user.Phone,
		&user.PasswordHash, &user.Role, &user.Department, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	query := `SELECT id, employee_code, full_name, email, phone, password_hash, role, department, status, created_at, updated_at FROM users WHERE id = $1`
	
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.EmployeeCode, &user.FullName, &user.Email, &user.Phone,
		&user.PasswordHash, &user.Role, &user.Department, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	query := `SELECT id, employee_code, full_name, email, phone, password_hash, role, department, status, created_at, updated_at FROM users`
	
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID, &user.EmployeeCode, &user.FullName, &user.Email, &user.Phone,
			&user.PasswordHash, &user.Role, &user.Department, &user.Status,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(user *model.User) error {
	query := `
		UPDATE users 
		SET employee_code = $1, full_name = $2, email = $3, phone = $4, password_hash = $5, role = $6, department = $7, status = $8, updated_at = NOW()
		WHERE id = $9`
	
	commandTag, err := r.db.Exec(context.Background(), query,
		user.EmployeeCode, user.FullName, user.Email, user.Phone,
		user.PasswordHash, user.Role, user.Department, user.Status, user.ID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *UserRepository) Search(searchQuery string) ([]model.User, error) {
	var users []model.User
	q := "%" + searchQuery + "%"
	
	query := `
		SELECT id, employee_code, full_name, email, phone, password_hash, role, department, status, created_at, updated_at 
		FROM users 
		WHERE employee_code ILIKE $1 OR full_name ILIKE $2 OR phone ILIKE $3 OR email ILIKE $4`
	
	rows, err := r.db.Query(context.Background(), query, q, q, q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID, &user.EmployeeCode, &user.FullName, &user.Email, &user.Phone,
			&user.PasswordHash, &user.Role, &user.Department, &user.Status,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
