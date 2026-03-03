-- Chart of Accounts — Standar Indonesia (PSAK EMKM)
-- Jalankan setelah migration: psql -U postgres -d journal_entry -f db/seeds/chart_of_accounts.sql

INSERT INTO accounts (code, name, type, description) VALUES
-- === ASET (Assets) ===
('1100', 'Kas',                        'asset',     'Uang tunai di tangan'),
('1200', 'Bank',                       'asset',     'Saldo rekening bank'),
('1300', 'Piutang Usaha',              'asset',     'Tagihan dari pelanggan'),
('1400', 'Persediaan',                 'asset',     'Barang dagangan / bahan baku'),
('1500', 'Beban Dibayar di Muka',      'asset',     'Pembayaran di muka (sewa, asuransi)'),
('1600', 'Aset Tetap',                 'asset',     'Tanah, bangunan, kendaraan, peralatan'),
('1700', 'Akumulasi Penyusutan',       'asset',     'Penyusutan aset tetap (contra account)'),

-- === KEWAJIBAN (Liabilities) ===
('2100', 'Utang Usaha',               'liability', 'Utang kepada supplier'),
('2200', 'Utang Gaji',                'liability', 'Gaji yang belum dibayar'),
('2300', 'Utang Pajak',               'liability', 'Pajak yang belum disetor'),
('2400', 'Utang Bank',                'liability', 'Pinjaman dari bank'),

-- === EKUITAS (Equity) ===
('3100', 'Modal Pemilik',             'equity',    'Modal yang disetor pemilik'),
('3200', 'Laba Ditahan',              'equity',    'Akumulasi laba dari periode sebelumnya'),
('3300', 'Prive',                     'equity',    'Penarikan modal oleh pemilik'),

-- === PENDAPATAN (Revenue) ===
('4100', 'Pendapatan Penjualan',      'revenue',   'Pendapatan dari penjualan barang'),
('4200', 'Pendapatan Jasa',           'revenue',   'Pendapatan dari penyediaan jasa'),
('4900', 'Pendapatan Lain-lain',      'revenue',   'Pendapatan non-operasional'),

-- === HARGA POKOK (Cost of Goods Sold) ===
('5100', 'Harga Pokok Penjualan',     'cogs',      'Biaya langsung untuk barang yang dijual'),

-- === BEBAN (Expenses) ===
('6100', 'Beban Gaji',                'expense',   'Gaji & tunjangan karyawan'),
('6200', 'Beban Sewa',                'expense',   'Biaya sewa tempat usaha'),
('6300', 'Beban Listrik, Air & Telepon', 'expense', 'Biaya utilitas'),
('6400', 'Beban Pemasaran',           'expense',   'Biaya iklan, promosi'),
('6500', 'Beban Perlengkapan',        'expense',   'Biaya ATK, supplies'),
('6600', 'Beban Penyusutan',          'expense',   'Penyusutan aset tetap per periode'),
('6700', 'Beban Transportasi',        'expense',   'Biaya perjalanan, bensin'),
('6800', 'Beban Administrasi',        'expense',   'Biaya admin & umum'),
('6900', 'Beban Lain-lain',           'expense',   'Beban non-operasional lainnya')

ON CONFLICT (code) DO NOTHING;
