package entity

// Song is the playlist module's own view of a song, populated through the
// SongCatalog port. It is deliberately independent of the song module's types.
type Song struct {
	ID              uint64
	Name            string
	ArtistID        uint64
	DurationSeconds uint32
	StorageLocation string
}
