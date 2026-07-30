package entity

// Playlist is the user module's own view of a playlist, populated through the
// PlaylistCatalog port. It is deliberately independent of the playlist module's types.
type Playlist struct {
	ID   uint64
	Name string
}
