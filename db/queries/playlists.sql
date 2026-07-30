-- name: CreatePlaylist :execresult
insert into playlists (name) values (?);

-- name: LinkPlaylistToUser :exec
insert into user_playlists (user_id, playlist_id) values (?, ?);

-- name: GetPlaylistByID :one
select * from playlists where id = ?;

-- name: AddSongToPlaylist :exec
insert ignore into playlist_songs (playlist_id, song_id) values (?, ?);

-- name: GetPlaylistSongIDs :many
select song_id from playlist_songs where playlist_id = ?;

-- name: GetPlaylistsByIDs :many
select * from playlists where id in (sqlc.slice('ids'));

-- name: GetPlaylistsByUserID :many
select p.* from playlists p
join user_playlists up on up.playlist_id = p.id
where up.user_id = ?;
