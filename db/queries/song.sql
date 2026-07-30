-- name: GetSongsByIDs :many
select * from songs where id in (sqlc.slice('ids'));
