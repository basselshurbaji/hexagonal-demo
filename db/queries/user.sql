-- name: GetUserByID :one
select * from users where id = ?;

-- name: GetUserPlaylistIDs :many
select playlist_id from user_playlists where user_id = ?;