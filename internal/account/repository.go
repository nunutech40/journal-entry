package account

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for account data access.
type Repository interface {
	GetAll(ctx context.Context) ([]Account, error)
	GetByID(ctx context.Context, id int) (*Account, error)
	GetByCode(ctx context.Context, code string) (*Account, error)
	Create(ctx context.Context, req CreateRequest) (*Account, error)
	Update(ctx context.Context, id int, req UpdateRequest) (*Account, error)
	Delete(ctx context.Context, id int) error
}

// pgxRepository implements Repository using pgxpool.
type pgxRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new pgx-based Repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) GetAll(ctx context.Context) ([]Account, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, code, name, type, parent_id, description, is_active, created_at, updated_at
		 FROM accounts
		 ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var a Account
		err := rows.Scan(&a.ID, &a.Code, &a.Name, &a.Type, &a.ParentID,
			&a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, a)
	}

	return accounts, rows.Err()
}

func (r *pgxRepository) GetByID(ctx context.Context, id int) (*Account, error) {
	var a Account
	err := r.pool.QueryRow(ctx,
		`SELECT id, code, name, type, parent_id, description, is_active, created_at, updated_at
		 FROM accounts WHERE id = $1`, id).
		Scan(&a.ID, &a.Code, &a.Name, &a.Type, &a.ParentID,
			&a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get account %d: %w", id, err)
	}
	return &a, nil
}

func (r *pgxRepository) GetByCode(ctx context.Context, code string) (*Account, error) {
	var a Account
	err := r.pool.QueryRow(ctx,
		`SELECT id, code, name, type, parent_id, description, is_active, created_at, updated_at
		 FROM accounts WHERE code = $1`, code).
		Scan(&a.ID, &a.Code, &a.Name, &a.Type, &a.ParentID,
			&a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get account by code %s: %w", code, err)
	}
	return &a, nil
}

func (r *pgxRepository) Create(ctx context.Context, req CreateRequest) (*Account, error) {
	var a Account
	err := r.pool.QueryRow(ctx,
		`INSERT INTO accounts (code, name, type, description)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, code, name, type, parent_id, description, is_active, created_at, updated_at`,
		req.Code, req.Name, req.Type, req.Description).
		Scan(&a.ID, &a.Code, &a.Name, &a.Type, &a.ParentID,
			&a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	return &a, nil
}

func (r *pgxRepository) Update(ctx context.Context, id int, req UpdateRequest) (*Account, error) {
	var a Account
	err := r.pool.QueryRow(ctx,
		`UPDATE accounts
		 SET code = $2, name = $3, type = $4, description = $5, is_active = $6, updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, code, name, type, parent_id, description, is_active, created_at, updated_at`,
		id, req.Code, req.Name, req.Type, req.Description, req.IsActive).
		Scan(&a.ID, &a.Code, &a.Name, &a.Type, &a.ParentID,
			&a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("update account %d: %w", id, err)
	}
	return &a, nil
}

func (r *pgxRepository) Delete(ctx context.Context, id int) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM accounts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete account %d: %w", id, err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("account %d not found", id)
	}
	return nil
}
