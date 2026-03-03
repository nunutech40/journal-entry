package report

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for report data access.
type Repository interface {
	GetLedger(ctx context.Context, accountID int, dateFrom, dateTo time.Time) ([]LedgerEntry, error)
	GetTrialBalance(ctx context.Context) ([]TrialBalanceRow, error)
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new pgx-based report Repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) GetLedger(ctx context.Context, accountID int, dateFrom, dateTo time.Time) ([]LedgerEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT je.entry_date, je.id, je.reference_number,
		        COALESCE(el.description, je.description) AS description,
		        el.debit, el.credit
		 FROM entry_lines el
		 JOIN journal_entries je ON je.id = el.journal_entry_id
		 WHERE el.account_id = $1
		   AND je.entry_date >= $2
		   AND je.entry_date <= $3
		 ORDER BY je.entry_date, je.id`,
		accountID, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("query ledger: %w", err)
	}
	defer rows.Close()

	var entries []LedgerEntry
	for rows.Next() {
		var e LedgerEntry
		err := rows.Scan(&e.Date, &e.JournalEntryID, &e.ReferenceNumber,
			&e.Description, &e.Debit, &e.Credit)
		if err != nil {
			return nil, fmt.Errorf("scan ledger entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *pgxRepository) GetTrialBalance(ctx context.Context) ([]TrialBalanceRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT a.id, a.code, a.name, a.type,
		        COALESCE(SUM(el.debit), 0) AS debit_total,
		        COALESCE(SUM(el.credit), 0) AS credit_total
		 FROM accounts a
		 LEFT JOIN entry_lines el ON el.account_id = a.id
		 LEFT JOIN journal_entries je ON je.id = el.journal_entry_id
		 WHERE a.is_active = true
		 GROUP BY a.id, a.code, a.name, a.type
		 ORDER BY a.code`)
	if err != nil {
		return nil, fmt.Errorf("query trial balance: %w", err)
	}
	defer rows.Close()

	var result []TrialBalanceRow
	for rows.Next() {
		var row TrialBalanceRow
		err := rows.Scan(&row.AccountID, &row.AccountCode, &row.AccountName, &row.AccountType,
			&row.DebitTotal, &row.CreditTotal)
		if err != nil {
			return nil, fmt.Errorf("scan trial balance row: %w", err)
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
