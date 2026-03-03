package account

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Common errors
var (
	ErrCodeRequired    = errors.New("kode akun wajib diisi")
	ErrNameRequired    = errors.New("nama akun wajib diisi")
	ErrInvalidType     = errors.New("tipe akun tidak valid")
	ErrCodeExists      = errors.New("kode akun sudah digunakan")
	ErrAccountNotFound = errors.New("akun tidak ditemukan")
)

// Service defines business logic for account operations.
type Service interface {
	ListAccounts(ctx context.Context) ([]Account, error)
	GetAccount(ctx context.Context, id int) (*Account, error)
	CreateAccount(ctx context.Context, req CreateRequest) (*Account, error)
	UpdateAccount(ctx context.Context, id int, req UpdateRequest) (*Account, error)
	DeleteAccount(ctx context.Context, id int) error
}

type serviceImpl struct {
	repo Repository
}

// NewService creates a new account Service.
func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) ListAccounts(ctx context.Context) ([]Account, error) {
	return s.repo.GetAll(ctx)
}

func (s *serviceImpl) GetAccount(ctx context.Context, id int) (*Account, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, ErrAccountNotFound
	}
	return acc, nil
}

func (s *serviceImpl) CreateAccount(ctx context.Context, req CreateRequest) (*Account, error) {
	// Trim whitespace
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Type = strings.TrimSpace(strings.ToLower(req.Type))
	req.Description = strings.TrimSpace(req.Description)

	// Validate
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	// Check code uniqueness
	existing, err := s.repo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("check code uniqueness: %w", err)
	}
	if existing != nil {
		return nil, ErrCodeExists
	}

	return s.repo.Create(ctx, req)
}

func (s *serviceImpl) UpdateAccount(ctx context.Context, id int, req UpdateRequest) (*Account, error) {
	// Trim whitespace
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Type = strings.TrimSpace(strings.ToLower(req.Type))
	req.Description = strings.TrimSpace(req.Description)

	// Validate
	if err := validateUpdate(req); err != nil {
		return nil, err
	}

	// Check if account exists
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrAccountNotFound
	}

	// Check code uniqueness (only if code changed)
	if req.Code != current.Code {
		existing, err := s.repo.GetByCode(ctx, req.Code)
		if err != nil {
			return nil, fmt.Errorf("check code uniqueness: %w", err)
		}
		if existing != nil {
			return nil, ErrCodeExists
		}
	}

	return s.repo.Update(ctx, id, req)
}

func (s *serviceImpl) DeleteAccount(ctx context.Context, id int) error {
	// Check if account exists
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountNotFound
	}

	return s.repo.Delete(ctx, id)
}

// --- Validation helpers ---

func validateCreate(req CreateRequest) error {
	if req.Code == "" {
		return ErrCodeRequired
	}
	if req.Name == "" {
		return ErrNameRequired
	}
	if !isValidType(req.Type) {
		return ErrInvalidType
	}
	return nil
}

func validateUpdate(req UpdateRequest) error {
	if req.Code == "" {
		return ErrCodeRequired
	}
	if req.Name == "" {
		return ErrNameRequired
	}
	if !isValidType(req.Type) {
		return ErrInvalidType
	}
	return nil
}

func isValidType(t string) bool {
	for _, valid := range ValidTypes() {
		if t == valid {
			return true
		}
	}
	return false
}
