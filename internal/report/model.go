package report

import "time"

// LedgerEntry represents a single row in the general ledger report.
type LedgerEntry struct {
	Date            time.Time `json:"date"`
	JournalEntryID  int       `json:"journal_entry_id"`
	ReferenceNumber string    `json:"reference_number"`
	Description     string    `json:"description"`
	Debit           float64   `json:"debit"`
	Credit          float64   `json:"credit"`
	Balance         float64   `json:"balance"` // running balance
}

// LedgerReport holds the full general ledger report for a single account.
type LedgerReport struct {
	AccountID   int           `json:"account_id"`
	AccountCode string        `json:"account_code"`
	AccountName string        `json:"account_name"`
	AccountType string        `json:"account_type"`
	DateFrom    string        `json:"date_from"`
	DateTo      string        `json:"date_to"`
	Entries     []LedgerEntry `json:"entries"`
	TotalDebit  float64       `json:"total_debit"`
	TotalCredit float64       `json:"total_credit"`
	EndBalance  float64       `json:"end_balance"`
}

// TrialBalanceRow represents a single account row in the trial balance.
type TrialBalanceRow struct {
	AccountID   int     `json:"account_id"`
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	AccountType string  `json:"account_type"`
	DebitTotal  float64 `json:"debit_total"`
	CreditTotal float64 `json:"credit_total"`
}

// TrialBalance holds the full trial balance report.
type TrialBalance struct {
	Rows        []TrialBalanceRow `json:"rows"`
	TotalDebit  float64           `json:"total_debit"`
	TotalCredit float64           `json:"total_credit"`
	IsBalanced  bool              `json:"is_balanced"`
}
