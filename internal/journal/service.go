package journal

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"journal-entry/internal/account"
)

// Common errors
var (
	ErrDateRequired        = errors.New("tanggal wajib diisi")
	ErrDateInvalid         = errors.New("format tanggal tidak valid (gunakan YYYY-MM-DD)")
	ErrDescriptionRequired = errors.New("keterangan wajib diisi")
	ErrMinTwoLines         = errors.New("minimal 2 baris entri diperlukan")
	ErrLineNoAccount       = errors.New("setiap baris harus memiliki akun")
	ErrLineNoAmount        = errors.New("setiap baris harus memiliki nominal debit atau kredit")
	ErrLineBothAmounts     = errors.New("setiap baris hanya boleh debit ATAU kredit, tidak keduanya")
	ErrNotBalanced         = errors.New("total debit harus sama dengan total kredit")
	ErrAccountNotFound     = errors.New("akun tidak ditemukan")
	ErrEntryNotFound       = errors.New("jurnal tidak ditemukan")
)

// Service defines business logic for journal entry operations.
type Service interface {
	ListEntries(ctx context.Context) ([]JournalEntry, error)
	GetEntry(ctx context.Context, id int) (*JournalEntry, error)
	CreateEntry(ctx context.Context, req CreateRequest) (*JournalEntry, error)
	UpdateEntry(ctx context.Context, id int, req UpdateRequest) (*JournalEntry, error)
	DeleteEntry(ctx context.Context, id int) error
}

type serviceImpl struct {
	repo        Repository
	accountRepo account.Repository
}

// NewService creates a new journal Service.
func NewService(repo Repository, accountRepo account.Repository) Service {
	return &serviceImpl{repo: repo, accountRepo: accountRepo}
}

func (s *serviceImpl) ListEntries(ctx context.Context) ([]JournalEntry, error) {
	return s.repo.GetAll(ctx)
}

func (s *serviceImpl) GetEntry(ctx context.Context, id int) (*JournalEntry, error) {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrEntryNotFound
	}
	return entry, nil
}

func (s *serviceImpl) CreateEntry(ctx context.Context, req CreateRequest) (*JournalEntry, error) {
	// Trim
	req.Description = strings.TrimSpace(req.Description)
	req.ReferenceNumber = strings.TrimSpace(req.ReferenceNumber)

	// Validate
	entryDate, err := s.validateRequest(ctx, req.EntryDate, req.Description, req.Lines)
	if err != nil {
		return nil, err
	}

	entry := &JournalEntry{
		EntryDate:       entryDate,
		Description:     req.Description,
		ReferenceNumber: req.ReferenceNumber,
	}

	return s.repo.Create(ctx, entry, req.Lines)
}

func (s *serviceImpl) UpdateEntry(ctx context.Context, id int, req UpdateRequest) (*JournalEntry, error) {
	// Check exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrEntryNotFound
	}

	// Trim
	req.Description = strings.TrimSpace(req.Description)
	req.ReferenceNumber = strings.TrimSpace(req.ReferenceNumber)

	// Validate
	entryDate, err := s.validateRequest(ctx, req.EntryDate, req.Description, req.Lines)
	if err != nil {
		return nil, err
	}

	entry := &JournalEntry{
		EntryDate:       entryDate,
		Description:     req.Description,
		ReferenceNumber: req.ReferenceNumber,
	}

	return s.repo.Update(ctx, id, entry, req.Lines)
}

func (s *serviceImpl) DeleteEntry(ctx context.Context, id int) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrEntryNotFound
	}
	return s.repo.Delete(ctx, id)
}

// validateRequest validates common fields for create/update.
func (s *serviceImpl) validateRequest(ctx context.Context, dateStr, description string, lines []CreateLineRequest) (time.Time, error) {
	// Date
	if strings.TrimSpace(dateStr) == "" {
		return time.Time{}, ErrDateRequired
	}
	entryDate, err := parseDate(strings.TrimSpace(dateStr))
	if err != nil {
		return time.Time{}, ErrDateInvalid
	}

	// Description
	if description == "" {
		return time.Time{}, ErrDescriptionRequired
	}

	// Lines
	if len(lines) < 2 {
		return time.Time{}, ErrMinTwoLines
	}

	var totalDebit, totalCredit float64
	for i, line := range lines {
		if line.AccountID == 0 {
			return time.Time{}, fmt.Errorf("baris %d: %w", i+1, ErrLineNoAccount)
		}
		if line.Debit == 0 && line.Credit == 0 {
			return time.Time{}, fmt.Errorf("baris %d: %w", i+1, ErrLineNoAmount)
		}
		if line.Debit > 0 && line.Credit > 0 {
			return time.Time{}, fmt.Errorf("baris %d: %w", i+1, ErrLineBothAmounts)
		}

		// Validate account exists
		acc, err := s.accountRepo.GetByID(ctx, line.AccountID)
		if err != nil {
			return time.Time{}, fmt.Errorf("cek akun baris %d: %w", i+1, err)
		}
		if acc == nil {
			return time.Time{}, fmt.Errorf("baris %d: %w (ID: %d)", i+1, ErrAccountNotFound, line.AccountID)
		}

		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	// Balance check (float comparison with tolerance)
	if math.Abs(totalDebit-totalCredit) > 0.001 {
		return time.Time{}, ErrNotBalanced
	}

	return entryDate, nil
}
