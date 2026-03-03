package journal

import (
	"context"
	"errors"
	"testing"
	"time"

	"journal-entry/internal/account"
)

// --- Mock Account Repository ---

type mockAccountRepo struct {
	accounts map[int]*account.Account
}

func newMockAccountRepo() *mockAccountRepo {
	m := &mockAccountRepo{accounts: make(map[int]*account.Account)}
	// Seed some accounts
	m.accounts[1] = &account.Account{ID: 1, Code: "1100", Name: "Kas", Type: "asset", IsActive: true}
	m.accounts[2] = &account.Account{ID: 2, Code: "4100", Name: "Pendapatan Penjualan", Type: "revenue", IsActive: true}
	m.accounts[3] = &account.Account{ID: 3, Code: "6100", Name: "Beban Gaji", Type: "expense", IsActive: true}
	return m
}

func (m *mockAccountRepo) GetAll(_ context.Context) ([]account.Account, error) { return nil, nil }
func (m *mockAccountRepo) GetByCode(_ context.Context, _ string) (*account.Account, error) {
	return nil, nil
}
func (m *mockAccountRepo) Create(_ context.Context, _ account.CreateRequest) (*account.Account, error) {
	return nil, nil
}
func (m *mockAccountRepo) Update(_ context.Context, _ int, _ account.UpdateRequest) (*account.Account, error) {
	return nil, nil
}
func (m *mockAccountRepo) Delete(_ context.Context, _ int) error { return nil }

func (m *mockAccountRepo) GetByID(_ context.Context, id int) (*account.Account, error) {
	a, ok := m.accounts[id]
	if !ok {
		return nil, nil
	}
	return a, nil
}

// --- Mock Journal Repository ---

type mockJournalRepo struct {
	entries map[int]*JournalEntry
	nextID  int
}

func newMockJournalRepo() *mockJournalRepo {
	return &mockJournalRepo{
		entries: make(map[int]*JournalEntry),
		nextID:  1,
	}
}

func (m *mockJournalRepo) GetAll(_ context.Context) ([]JournalEntry, error) {
	var result []JournalEntry
	for _, e := range m.entries {
		result = append(result, *e)
	}
	return result, nil
}

func (m *mockJournalRepo) GetByID(_ context.Context, id int) (*JournalEntry, error) {
	e, ok := m.entries[id]
	if !ok {
		return nil, nil
	}
	return e, nil
}

func (m *mockJournalRepo) Create(_ context.Context, entry *JournalEntry, lines []CreateLineRequest) (*JournalEntry, error) {
	entry.ID = m.nextID
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	m.entries[entry.ID] = entry
	m.nextID++
	return entry, nil
}

func (m *mockJournalRepo) Update(_ context.Context, id int, entry *JournalEntry, lines []CreateLineRequest) (*JournalEntry, error) {
	if _, ok := m.entries[id]; !ok {
		return nil, nil
	}
	entry.ID = id
	entry.UpdatedAt = time.Now()
	m.entries[id] = entry
	return entry, nil
}

func (m *mockJournalRepo) Delete(_ context.Context, id int) error {
	if _, ok := m.entries[id]; !ok {
		return errors.New("not found")
	}
	delete(m.entries, id)
	return nil
}

// --- Tests ---

func validLines() []CreateLineRequest {
	return []CreateLineRequest{
		{AccountID: 1, Debit: 100000, Credit: 0},
		{AccountID: 2, Debit: 0, Credit: 100000},
	}
}

func TestCreateEntry_Valid(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	entry, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "2026-01-15",
		Description: "Penjualan tunai",
		Lines:       validLines(),
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if entry.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if entry.Description != "Penjualan tunai" {
		t.Errorf("expected description 'Penjualan tunai', got '%s'", entry.Description)
	}
}

func TestCreateEntry_EmptyDate(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	_, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "",
		Description: "Test",
		Lines:       validLines(),
	})
	if !errors.Is(err, ErrDateRequired) {
		t.Errorf("expected ErrDateRequired, got: %v", err)
	}
}

func TestCreateEntry_InvalidDate(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	_, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "31-01-2026",
		Description: "Test",
		Lines:       validLines(),
	})
	if !errors.Is(err, ErrDateInvalid) {
		t.Errorf("expected ErrDateInvalid, got: %v", err)
	}
}

func TestCreateEntry_EmptyDescription(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	_, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "2026-01-15",
		Description: "",
		Lines:       validLines(),
	})
	if !errors.Is(err, ErrDescriptionRequired) {
		t.Errorf("expected ErrDescriptionRequired, got: %v", err)
	}
}

func TestCreateEntry_SingleLine(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	_, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "2026-01-15",
		Description: "Test",
		Lines: []CreateLineRequest{
			{AccountID: 1, Debit: 100000, Credit: 0},
		},
	})
	if !errors.Is(err, ErrMinTwoLines) {
		t.Errorf("expected ErrMinTwoLines, got: %v", err)
	}
}

func TestCreateEntry_LineBothDebitAndCredit(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	_, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "2026-01-15",
		Description: "Test",
		Lines: []CreateLineRequest{
			{AccountID: 1, Debit: 100000, Credit: 50000},
			{AccountID: 2, Debit: 0, Credit: 50000},
		},
	})
	if err == nil {
		t.Fatal("expected error for line with both debit and credit")
	}
	if !errors.Is(err, ErrLineBothAmounts) {
		t.Errorf("expected ErrLineBothAmounts, got: %v", err)
	}
}

func TestCreateEntry_NotBalanced(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	_, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "2026-01-15",
		Description: "Test",
		Lines: []CreateLineRequest{
			{AccountID: 1, Debit: 100000, Credit: 0},
			{AccountID: 2, Debit: 0, Credit: 50000},
		},
	})
	if !errors.Is(err, ErrNotBalanced) {
		t.Errorf("expected ErrNotBalanced, got: %v", err)
	}
}

func TestCreateEntry_AccountNotFound(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	_, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "2026-01-15",
		Description: "Test",
		Lines: []CreateLineRequest{
			{AccountID: 999, Debit: 100000, Credit: 0},
			{AccountID: 2, Debit: 0, Credit: 100000},
		},
	})
	if err == nil {
		t.Fatal("expected error for non-existent account")
	}
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got: %v", err)
	}
}

func TestCreateEntry_LineNoAmount(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	_, err := svc.CreateEntry(context.Background(), CreateRequest{
		EntryDate:   "2026-01-15",
		Description: "Test",
		Lines: []CreateLineRequest{
			{AccountID: 1, Debit: 0, Credit: 0},
			{AccountID: 2, Debit: 0, Credit: 100000},
		},
	})
	if !errors.Is(err, ErrLineNoAmount) {
		t.Errorf("expected ErrLineNoAmount, got: %v", err)
	}
}

func TestDeleteEntry_NotFound(t *testing.T) {
	svc := NewService(newMockJournalRepo(), newMockAccountRepo())

	err := svc.DeleteEntry(context.Background(), 999)
	if !errors.Is(err, ErrEntryNotFound) {
		t.Errorf("expected ErrEntryNotFound, got: %v", err)
	}
}
