-- +goose Up
CREATE TABLE users (
    id    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name  VARCHAR(255)    NOT NULL,
    email VARCHAR(255)    NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_users_email (email)
);

CREATE TABLE artists (
    id   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255)    NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE songs (
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name             VARCHAR(255)    NOT NULL,
    artist_id        BIGINT UNSIGNED NOT NULL,
    duration_seconds INT UNSIGNED    NOT NULL,
    storage_location VARCHAR(512)    NOT NULL,
    PRIMARY KEY (id),
    KEY idx_songs_artist_id (artist_id),
    CONSTRAINT fk_songs_artist FOREIGN KEY (artist_id) REFERENCES artists (id)
);

CREATE TABLE playlists (
    id   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255)    NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE playlist_songs (
    playlist_id BIGINT UNSIGNED NOT NULL,
    song_id     BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (playlist_id, song_id),
    KEY idx_playlist_songs_song_id (song_id),
    CONSTRAINT fk_playlist_songs_playlist FOREIGN KEY (playlist_id) REFERENCES playlists (id) ON DELETE CASCADE,
    CONSTRAINT fk_playlist_songs_song FOREIGN KEY (song_id) REFERENCES songs (id) ON DELETE CASCADE
);

CREATE TABLE user_playlists (
    user_id     BIGINT UNSIGNED NOT NULL,
    playlist_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, playlist_id),
    KEY idx_user_playlists_playlist_id (playlist_id),
    CONSTRAINT fk_user_playlists_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_playlists_playlist FOREIGN KEY (playlist_id) REFERENCES playlists (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE user_playlists;
DROP TABLE playlist_songs;
DROP TABLE playlists;
DROP TABLE songs;
DROP TABLE artists;
DROP TABLE users;
