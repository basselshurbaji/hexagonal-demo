package adapters

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"hexagonal-demo/modules/user"
)

// HTTPAdapter is a driving adapter: it talks to the module through its public
// facade, like any other consumer.
type HTTPAdapter struct {
	users *user.Module
}

func NewHTTPAdapter(users *user.Module) *HTTPAdapter {
	return &HTTPAdapter{users: users}
}

func (h *HTTPAdapter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/{id}", h.getUserByID)
	mux.HandleFunc("GET /users/{id}/playlists", h.getUserPlaylists)
}

// userResponse is the transport representation — the model is not the wire format.
type userResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type playlistResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (h *HTTPAdapter) getUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	u, err := h.users.GetUserByID(r.Context(), id)
	switch {
	case errors.Is(err, user.ErrNotFound):
		http.Error(w, "user not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get user %d: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	})
}

func (h *HTTPAdapter) getUserPlaylists(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	playlists, err := h.users.GetUserPlaylists(r.Context(), id)
	switch {
	case errors.Is(err, user.ErrNotFound):
		http.Error(w, "user not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get playlists for user %d: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := make([]playlistResponse, 0, len(playlists))
	for _, p := range playlists {
		resp = append(resp, playlistResponse{ID: p.ID, Name: p.Name})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
