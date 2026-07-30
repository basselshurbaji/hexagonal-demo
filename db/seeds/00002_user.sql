-- Seed data: one user, no playlists.

-- +goose Up
INSERT INTO users (id, name, email) VALUES
    (1, 'Demo User', 'demo@hexagonal.local');

-- +goose Down
DELETE FROM users;
