-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_transactions_waktu ON transactions(waktu);
CREATE INDEX idx_transactions_jenis ON transactions(jenis);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_transactions_waktu;
DROP INDEX idx_transactions_jenis;
-- +goose StatementEnd