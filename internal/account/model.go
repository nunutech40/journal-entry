package account

import "time"

// Account represents a single account in the Chart of Accounts.
type Account struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // asset, liability, equity, revenue, cogs, expense
	ParentID    *int      `json:"parent_id,omitempty"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateRequest holds data needed to create a new account.
type CreateRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// UpdateRequest holds data needed to update an existing account.
type UpdateRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// ValidTypes returns the list of valid account types.
func ValidTypes() []string {
	return []string{"asset", "liability", "equity", "revenue", "cogs", "expense"}
}

// TypeLabel returns the Indonesian label for an account type.
func TypeLabel(t string) string {
	labels := map[string]string{
		"asset":     "Aset",
		"liability": "Kewajiban",
		"equity":    "Ekuitas",
		"revenue":   "Pendapatan",
		"cogs":      "Harga Pokok",
		"expense":   "Beban",
	}
	if label, ok := labels[t]; ok {
		return label
	}
	return t
}
