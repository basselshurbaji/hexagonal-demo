package adapters

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hexagonal-demo/modules/song"
)

// HTTPAdapter is a driving adapter: it talks to the module through its public
// facade, like any other consumer.
type HTTPAdapter struct {
	songs *song.Module
}

func NewHTTPAdapter(songs *song.Module) *HTTPAdapter {
	return &HTTPAdapter{songs: songs}
}

func (h *HTTPAdapter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /songs", h.getSongsByIDs)
}

// songResponse is the transport representation — the model is not the wire format.
type songResponse struct {
	ID              uint64 `json:"id"`
	Name            string `json:"name"`
	ArtistID        uint64 `json:"artist_id"`
	DurationSeconds uint32 `json:"duration_seconds"`
	StorageLocation string `json:"storage_location"`
}

// getSongsByIDs handles GET /songs?ids=1,2,3
func (h *HTTPAdapter) getSongsByIDs(w http.ResponseWriter, r *http.Request) {
	rawIDs := r.URL.Query().Get("ids")
	if rawIDs == "" {
		http.Error(w, "missing ids query parameter", http.StatusBadRequest)
		return
	}
	parts := strings.Split(rawIDs, ",")
	ids := make([]uint64, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil {
			http.Error(w, "invalid song id: "+part, http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	songs, err := h.songs.GetSongsByIDs(r.Context(), ids)
	if err != nil {
		log.Printf("get songs %v: %v", ids, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := make([]songResponse, 0, len(songs))
	for _, s := range songs {
		resp = append(resp, songResponse{
			ID:              s.ID,
			Name:            s.Name,
			ArtistID:        s.ArtistID,
			DurationSeconds: s.DurationSeconds,
			StorageLocation: s.StorageLocation,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
