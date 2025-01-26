-- +goose Up
-- +goose StatementBegin
CREATE TYPE transaction_type AS ENUM ('tabung', 'tarik');

CREATE TABLE accounts
(
    no_rekening VARCHAR(20) PRIMARY KEY,
    nama        VARCHAR(255)       NOT NULL,
    nik         VARCHAR(16) UNIQUE NOT NULL,
    no_hp       VARCHAR(15) UNIQUE NOT NULL,
    saldo       BIGINT             NOT NULL DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE    DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions
(
    id          SERIAL PRIMARY KEY,
    no_rekening VARCHAR(20)      NOT NULL REFERENCES accounts (no_rekening),
    jenis       transaction_type NOT NULL,
    nominal     BIGINT           NOT NULL,
    waktu       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_accounts_nik ON accounts (nik);
CREATE INDEX idx_accounts_no_hp ON accounts (no_hp);
CREATE INDEX idx_transactions_no_rekening ON transactions (no_rekening);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
DROP TABLE accounts;
DROP TYPE transaction_type;
-- +goose StatementEnd