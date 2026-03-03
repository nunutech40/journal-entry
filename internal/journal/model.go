package journal

import "time"

// JournalEntry represents a journal entry header.
type JournalEntry struct {
	ID              int         `json:"id"`
	EntryDate       time.Time   `json:"entry_date"`
	Description     string      `json:"description"`
	ReferenceNumber string      `json:"reference_number"`
	IsPosted        bool        `json:"is_posted"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Lines           []EntryLine `json:"lines"`
	TotalDebit      float64     `json:"total_debit"`
	TotalCredit     float64     `json:"total_credit"`
}

// EntryLine represents a single debit/credit line in a journal entry.
type EntryLine struct {
	ID             int     `json:"id"`
	JournalEntryID int     `json:"journal_entry_id"`
	AccountID      int     `json:"account_id"`
	AccountCode    string  `json:"account_code"` // from JOIN
	AccountName    string  `json:"account_name"` // from JOIN
	Description    string  `json:"description"`
	Debit          float64 `json:"debit"`
	Credit         float64 `json:"credit"`
}

// CreateRequest holds data for creating a journal entry.
type CreateRequest struct {
	EntryDate       string              `json:"entry_date"`
	Description     string              `json:"description"`
	ReferenceNumber string              `json:"reference_number"`
	Lines           []CreateLineRequest `json:"lines"`
}

// CreateLineRequest holds data for a single entry line.
type CreateLineRequest struct {
	AccountID   int     `json:"account_id"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}

// UpdateRequest holds data for updating a journal entry.
type UpdateRequest struct {
	EntryDate       string              `json:"entry_date"`
	Description     string              `json:"description"`
	ReferenceNumber string              `json:"reference_number"`
	Lines           []CreateLineRequest `json:"lines"`
}

// FormatDate returns the entry date as yyyy-mm-dd string for form input.
func (j *JournalEntry) FormatDate() string {
	return j.EntryDate.Format("2006-01-02")
}
