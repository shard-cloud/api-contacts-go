-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_contacts_deleted_at;
DROP INDEX IF EXISTS idx_contacts_email;
DROP TABLE IF EXISTS contacts;
-- +goose StatementEnd
