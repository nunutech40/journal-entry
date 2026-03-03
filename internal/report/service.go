package report

import (
	"context"
	"math"
	"time"

	"journal-entry/internal/account"
)

// Service defines business logic for reports.
type Service interface {
	GetLedgerReport(ctx context.Context, accountID int, dateFrom, dateTo string) (*LedgerReport, error)
	GetTrialBalance(ctx context.Context) (*TrialBalance, error)
}

type serviceImpl struct {
	repo        Repository
	accountRepo account.Repository
}

// NewService creates a new report Service.
func NewService(repo Repository, accountRepo account.Repository) Service {
	return &serviceImpl{repo: repo, accountRepo: accountRepo}
}

func (s *serviceImpl) GetLedgerReport(ctx context.Context, accountID int, dateFromStr, dateToStr string) (*LedgerReport, error) {
	// Get account info
	acc, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return &LedgerReport{}, nil
	}

	// Parse dates with defaults
	dateFrom, dateTo := parseDateRange(dateFromStr, dateToStr)

	// Get ledger entries
	entries, err := s.repo.GetLedger(ctx, accountID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Calculate running balance
	report := &LedgerReport{
		AccountID:   acc.ID,
		AccountCode: acc.Code,
		AccountName: acc.Name,
		AccountType: acc.Type,
		DateFrom:    dateFrom.Format("2006-01-02"),
		DateTo:      dateTo.Format("2006-01-02"),
	}

	var balance float64
	for i, e := range entries {
		// For asset/expense accounts: debit increases, credit decreases
		// For liability/equity/revenue: credit increases, debit decreases
		if isDebitNormal(acc.Type) {
			balance += e.Debit - e.Credit
		} else {
			balance += e.Credit - e.Debit
		}
		entries[i].Balance = balance
		report.TotalDebit += e.Debit
		report.TotalCredit += e.Credit
	}

	report.Entries = entries
	report.EndBalance = balance

	return report, nil
}

func (s *serviceImpl) GetTrialBalance(ctx context.Context) (*TrialBalance, error) {
	rows, err := s.repo.GetTrialBalance(ctx)
	if err != nil {
		return nil, err
	}

	tb := &TrialBalance{}
	for _, row := range rows {
		tb.Rows = append(tb.Rows, row)
		tb.TotalDebit += row.DebitTotal
		tb.TotalCredit += row.CreditTotal
	}

	tb.IsBalanced = math.Abs(tb.TotalDebit-tb.TotalCredit) < 0.01

	return tb, nil
}

// isDebitNormal returns true for account types that normally have debit balances.
func isDebitNormal(accountType string) bool {
	switch accountType {
	case "asset", "expense", "cogs":
		return true
	default:
		return false
	}
}

// parseDateRange parses date range strings, defaulting to start of year → today.
func parseDateRange(from, to string) (time.Time, time.Time) {
	now := time.Now()
	dateFrom := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	dateTo := now

	if from != "" {
		if parsed, err := time.Parse("2006-01-02", from); err == nil {
			dateFrom = parsed
		}
	}
	if to != "" {
		if parsed, err := time.Parse("2006-01-02", to); err == nil {
			dateTo = parsed
		}
	}

	return dateFrom, dateTo
}
