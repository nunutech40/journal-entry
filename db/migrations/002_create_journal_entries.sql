-- +goose Up
CREATE TABLE journal_entries (
    id SERIAL PRIMARY KEY,
    entry_date DATE NOT NULL,
    description TEXT NOT NULL,
    reference_number VARCHAR(50),
    is_posted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_journal_entries_date ON journal_entries(entry_date);
CREATE INDEX idx_journal_entries_posted ON journal_entries(is_posted);

-- +goose Down
DROP TABLE IF EXISTS journal_entries;
