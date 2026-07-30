package entity

type Song struct {
	ID              uint64
	Name            string
	ArtistID        uint64
	DurationSeconds uint32
	StorageLocation string
}
