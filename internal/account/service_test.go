package account

import (
	"context"
	"errors"
	"testing"
)

// --- Mock Repository ---

type mockRepo struct {
	accounts  map[int]*Account
	byCode    map[string]*Account
	nextID    int
	createErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		accounts: make(map[int]*Account),
		byCode:   make(map[string]*Account),
		nextID:   1,
	}
}

func (m *mockRepo) GetAll(_ context.Context) ([]Account, error) {
	var result []Account
	for _, a := range m.accounts {
		result = append(result, *a)
	}
	return result, nil
}

func (m *mockRepo) GetByID(_ context.Context, id int) (*Account, error) {
	a, ok := m.accounts[id]
	if !ok {
		return nil, nil
	}
	return a, nil
}

func (m *mockRepo) GetByCode(_ context.Context, code string) (*Account, error) {
	a, ok := m.byCode[code]
	if !ok {
		return nil, nil
	}
	return a, nil
}

func (m *mockRepo) Create(_ context.Context, req CreateRequest) (*Account, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	a := &Account{
		ID:       m.nextID,
		Code:     req.Code,
		Name:     req.Name,
		Type:     req.Type,
		IsActive: true,
	}
	m.accounts[a.ID] = a
	m.byCode[a.Code] = a
	m.nextID++
	return a, nil
}

func (m *mockRepo) Update(_ context.Context, id int, req UpdateRequest) (*Account, error) {
	a, ok := m.accounts[id]
	if !ok {
		return nil, nil
	}
	delete(m.byCode, a.Code)
	a.Code = req.Code
	a.Name = req.Name
	a.Type = req.Type
	a.Description = req.Description
	a.IsActive = req.IsActive
	m.byCode[a.Code] = a
	return a, nil
}

func (m *mockRepo) Delete(_ context.Context, id int) error {
	a, ok := m.accounts[id]
	if !ok {
		return errors.New("not found")
	}
	delete(m.byCode, a.Code)
	delete(m.accounts, id)
	return nil
}

// Helper: seed an account into mock repo
func (m *mockRepo) seed(code, name, typ string) *Account {
	a := &Account{ID: m.nextID, Code: code, Name: name, Type: typ, IsActive: true}
	m.accounts[a.ID] = a
	m.byCode[a.Code] = a
	m.nextID++
	return a
}

// --- Tests ---

func TestCreateAccount_Valid(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	acc, err := svc.CreateAccount(context.Background(), CreateRequest{
		Code: "1100",
		Name: "Kas",
		Type: "asset",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if acc.Code != "1100" {
		t.Errorf("expected code 1100, got %s", acc.Code)
	}
	if acc.Name != "Kas" {
		t.Errorf("expected name Kas, got %s", acc.Name)
	}
}

func TestCreateAccount_EmptyCode(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateAccount(context.Background(), CreateRequest{
		Code: "",
		Name: "Kas",
		Type: "asset",
	})
	if !errors.Is(err, ErrCodeRequired) {
		t.Errorf("expected ErrCodeRequired, got: %v", err)
	}
}

func TestCreateAccount_EmptyName(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateAccount(context.Background(), CreateRequest{
		Code: "1100",
		Name: "",
		Type: "asset",
	})
	if !errors.Is(err, ErrNameRequired) {
		t.Errorf("expected ErrNameRequired, got: %v", err)
	}
}

func TestCreateAccount_InvalidType(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateAccount(context.Background(), CreateRequest{
		Code: "1100",
		Name: "Kas",
		Type: "invalid",
	})
	if !errors.Is(err, ErrInvalidType) {
		t.Errorf("expected ErrInvalidType, got: %v", err)
	}
}

func TestCreateAccount_DuplicateCode(t *testing.T) {
	repo := newMockRepo()
	repo.seed("1100", "Kas", "asset")
	svc := NewService(repo)

	_, err := svc.CreateAccount(context.Background(), CreateRequest{
		Code: "1100",
		Name: "Kas Duplikat",
		Type: "asset",
	})
	if !errors.Is(err, ErrCodeExists) {
		t.Errorf("expected ErrCodeExists, got: %v", err)
	}
}

func TestCreateAccount_TrimsWhitespace(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	acc, err := svc.CreateAccount(context.Background(), CreateRequest{
		Code: "  1100  ",
		Name: "  Kas  ",
		Type: "  ASSET  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if acc.Code != "1100" {
		t.Errorf("expected trimmed code '1100', got '%s'", acc.Code)
	}
	if acc.Name != "Kas" {
		t.Errorf("expected trimmed name 'Kas', got '%s'", acc.Name)
	}
}

func TestUpdateAccount_Valid(t *testing.T) {
	repo := newMockRepo()
	existing := repo.seed("1100", "Kas", "asset")
	svc := NewService(repo)

	acc, err := svc.UpdateAccount(context.Background(), existing.ID, UpdateRequest{
		Code:     "1100",
		Name:     "Kas Besar",
		Type:     "asset",
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if acc.Name != "Kas Besar" {
		t.Errorf("expected name 'Kas Besar', got '%s'", acc.Name)
	}
}

func TestUpdateAccount_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.UpdateAccount(context.Background(), 999, UpdateRequest{
		Code: "1100",
		Name: "Kas",
		Type: "asset",
	})
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got: %v", err)
	}
}

func TestUpdateAccount_DuplicateCode(t *testing.T) {
	repo := newMockRepo()
	repo.seed("1100", "Kas", "asset")
	existing := repo.seed("1200", "Bank", "asset")
	svc := NewService(repo)

	_, err := svc.UpdateAccount(context.Background(), existing.ID, UpdateRequest{
		Code:     "1100", // already exists
		Name:     "Bank",
		Type:     "asset",
		IsActive: true,
	})
	if !errors.Is(err, ErrCodeExists) {
		t.Errorf("expected ErrCodeExists, got: %v", err)
	}
}

func TestDeleteAccount_Valid(t *testing.T) {
	repo := newMockRepo()
	existing := repo.seed("1100", "Kas", "asset")
	svc := NewService(repo)

	err := svc.DeleteAccount(context.Background(), existing.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDeleteAccount_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	err := svc.DeleteAccount(context.Background(), 999)
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got: %v", err)
	}
}
