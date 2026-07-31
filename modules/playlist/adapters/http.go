package adapters

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"hexagonal-demo/modules/playlist"
)

// HTTPAdapter is a driving adapter: it talks to the module through its public
// facade, like any other consumer.
type HTTPAdapter struct {
	playlists playlist.Facade
}

func NewHTTPAdapter(playlists playlist.Facade) *HTTPAdapter {
	return &HTTPAdapter{playlists: playlists}
}

func (h *HTTPAdapter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /playlists", h.createPlaylist)
	mux.HandleFunc("PUT /playlists/{id}/songs", h.addSongs)
	mux.HandleFunc("GET /playlists/{id}", h.getPlaylist)
}

type createPlaylistRequest struct {
	Name   string `json:"name"`
	UserID uint64 `json:"user_id"`
}

type playlistResponse struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	UserID uint64 `json:"user_id,omitempty"`
}

type addSongsRequest struct {
	SongIDs []uint64 `json:"song_ids"`
}

type songResponse struct {
	ID              uint64 `json:"id"`
	Name            string `json:"name"`
	ArtistID        uint64 `json:"artist_id"`
	DurationSeconds uint32 `json:"duration_seconds"`
	StorageLocation string `json:"storage_location"`
}

type playlistWithSongsResponse struct {
	ID    uint64         `json:"id"`
	Name  string         `json:"name"`
	Songs []songResponse `json:"songs"`
}

func (h *HTTPAdapter) createPlaylist(w http.ResponseWriter, r *http.Request) {
	var req createPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.UserID == 0 {
		http.Error(w, "name and user_id are required", http.StatusBadRequest)
		return
	}

	created, err := h.playlists.CreatePlaylist(r.Context(), req.Name, req.UserID)
	if err != nil {
		log.Printf("create playlist: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(playlistResponse{
		ID:     created.ID,
		Name:   created.Name,
		UserID: req.UserID,
	})
}

func (h *HTTPAdapter) addSongs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid playlist id", http.StatusBadRequest)
		return
	}
	var req addSongsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if len(req.SongIDs) == 0 {
		http.Error(w, "song_ids is required", http.StatusBadRequest)
		return
	}

	err = h.playlists.AddSongs(r.Context(), id, req.SongIDs)
	switch {
	case errors.Is(err, playlist.ErrNotFound):
		http.Error(w, "playlist not found", http.StatusNotFound)
		return
	case errors.Is(err, playlist.ErrSongNotFound):
		http.Error(w, "one or more songs do not exist", http.StatusBadRequest)
		return
	case err != nil:
		log.Printf("add songs to playlist %d: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPAdapter) getPlaylist(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid playlist id", http.StatusBadRequest)
		return
	}

	p, err := h.playlists.GetPlaylist(r.Context(), id)
	switch {
	case errors.Is(err, playlist.ErrNotFound):
		http.Error(w, "playlist not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get playlist %d: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := playlistWithSongsResponse{
		ID:    p.ID,
		Name:  p.Name,
		Songs: make([]songResponse, 0, len(p.Songs)),
	}
	for _, s := range p.Songs {
		resp.Songs = append(resp.Songs, songResponse{
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
