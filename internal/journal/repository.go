package journal

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for journal entry data access.
type Repository interface {
	GetAll(ctx context.Context) ([]JournalEntry, error)
	GetByID(ctx context.Context, id int) (*JournalEntry, error)
	Create(ctx context.Context, entry *JournalEntry, lines []CreateLineRequest) (*JournalEntry, error)
	Update(ctx context.Context, id int, entry *JournalEntry, lines []CreateLineRequest) (*JournalEntry, error)
	Delete(ctx context.Context, id int) error
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new pgx-based Repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) GetAll(ctx context.Context) ([]JournalEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT je.id, je.entry_date, je.description, je.reference_number,
		        je.is_posted, je.created_at, je.updated_at,
		        COALESCE(SUM(el.debit), 0) AS total_debit,
		        COALESCE(SUM(el.credit), 0) AS total_credit
		 FROM journal_entries je
		 LEFT JOIN entry_lines el ON el.journal_entry_id = je.id
		 GROUP BY je.id
		 ORDER BY je.entry_date DESC, je.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("query journal entries: %w", err)
	}
	defer rows.Close()

	var entries []JournalEntry
	for rows.Next() {
		var je JournalEntry
		err := rows.Scan(&je.ID, &je.EntryDate, &je.Description, &je.ReferenceNumber,
			&je.IsPosted, &je.CreatedAt, &je.UpdatedAt,
			&je.TotalDebit, &je.TotalCredit)
		if err != nil {
			return nil, fmt.Errorf("scan journal entry: %w", err)
		}
		entries = append(entries, je)
	}
	return entries, rows.Err()
}

func (r *pgxRepository) GetByID(ctx context.Context, id int) (*JournalEntry, error) {
	// Get header
	var je JournalEntry
	err := r.pool.QueryRow(ctx,
		`SELECT id, entry_date, description, reference_number, is_posted, created_at, updated_at
		 FROM journal_entries WHERE id = $1`, id).
		Scan(&je.ID, &je.EntryDate, &je.Description, &je.ReferenceNumber,
			&je.IsPosted, &je.CreatedAt, &je.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get journal entry %d: %w", id, err)
	}

	// Get lines with account info
	rows, err := r.pool.Query(ctx,
		`SELECT el.id, el.journal_entry_id, el.account_id,
		        a.code AS account_code, a.name AS account_name,
		        el.description, el.debit, el.credit
		 FROM entry_lines el
		 JOIN accounts a ON a.id = el.account_id
		 WHERE el.journal_entry_id = $1
		 ORDER BY el.id`, id)
	if err != nil {
		return nil, fmt.Errorf("query entry lines: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var line EntryLine
		err := rows.Scan(&line.ID, &line.JournalEntryID, &line.AccountID,
			&line.AccountCode, &line.AccountName,
			&line.Description, &line.Debit, &line.Credit)
		if err != nil {
			return nil, fmt.Errorf("scan entry line: %w", err)
		}
		je.Lines = append(je.Lines, line)
		je.TotalDebit += line.Debit
		je.TotalCredit += line.Credit
	}

	return &je, rows.Err()
}

func (r *pgxRepository) Create(ctx context.Context, entry *JournalEntry, lines []CreateLineRequest) (*JournalEntry, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert header
	var je JournalEntry
	err = tx.QueryRow(ctx,
		`INSERT INTO journal_entries (entry_date, description, reference_number)
		 VALUES ($1, $2, $3)
		 RETURNING id, entry_date, description, reference_number, is_posted, created_at, updated_at`,
		entry.EntryDate, entry.Description, entry.ReferenceNumber).
		Scan(&je.ID, &je.EntryDate, &je.Description, &je.ReferenceNumber,
			&je.IsPosted, &je.CreatedAt, &je.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert journal entry: %w", err)
	}

	// Insert lines
	for _, line := range lines {
		_, err := tx.Exec(ctx,
			`INSERT INTO entry_lines (journal_entry_id, account_id, description, debit, credit)
			 VALUES ($1, $2, $3, $4, $5)`,
			je.ID, line.AccountID, line.Description, line.Debit, line.Credit)
		if err != nil {
			return nil, fmt.Errorf("insert entry line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &je, nil
}

func (r *pgxRepository) Update(ctx context.Context, id int, entry *JournalEntry, lines []CreateLineRequest) (*JournalEntry, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update header
	var je JournalEntry
	err = tx.QueryRow(ctx,
		`UPDATE journal_entries
		 SET entry_date = $2, description = $3, reference_number = $4, updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, entry_date, description, reference_number, is_posted, created_at, updated_at`,
		id, entry.EntryDate, entry.Description, entry.ReferenceNumber).
		Scan(&je.ID, &je.EntryDate, &je.Description, &je.ReferenceNumber,
			&je.IsPosted, &je.CreatedAt, &je.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("update journal entry: %w", err)
	}

	// Delete old lines, insert new ones
	_, err = tx.Exec(ctx, `DELETE FROM entry_lines WHERE journal_entry_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("delete old lines: %w", err)
	}

	for _, line := range lines {
		_, err := tx.Exec(ctx,
			`INSERT INTO entry_lines (journal_entry_id, account_id, description, debit, credit)
			 VALUES ($1, $2, $3, $4, $5)`,
			je.ID, line.AccountID, line.Description, line.Debit, line.Credit)
		if err != nil {
			return nil, fmt.Errorf("insert entry line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &je, nil
}

func (r *pgxRepository) Delete(ctx context.Context, id int) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM journal_entries WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete journal entry %d: %w", id, err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("journal entry %d not found", id)
	}
	return nil
}

// parseDate parses a date string in yyyy-mm-dd format.
func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
