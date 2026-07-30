-- Bootstrap no-op migration: exists only so `make migrate` can be validated
-- before any real domain migrations are written. Safe to delete once you add
-- your first real migration (delete before running migrate, or just leave it).

-- +goose Up
SELECT 1;

-- +goose Down
SELECT 1;