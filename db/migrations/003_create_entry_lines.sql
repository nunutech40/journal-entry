-- +goose Up
CREATE TABLE entry_lines (
    id SERIAL PRIMARY KEY,
    journal_entry_id INTEGER NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    description TEXT,
    debit NUMERIC(15,2) DEFAULT 0 CHECK (debit >= 0),
    credit NUMERIC(15,2) DEFAULT 0 CHECK (credit >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),

    -- Setiap line harus debit ATAU credit, nggak boleh dua-duanya > 0
    CONSTRAINT chk_debit_xor_credit CHECK (
        (debit > 0 AND credit = 0) OR (debit = 0 AND credit > 0)
    )
);

CREATE INDEX idx_entry_lines_journal ON entry_lines(journal_entry_id);
CREATE INDEX idx_entry_lines_account ON entry_lines(account_id);

-- +goose Down
DROP TABLE IF EXISTS entry_lines;
