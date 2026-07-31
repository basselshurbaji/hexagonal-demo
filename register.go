package main

import (
	"database/sql"
	"net/http"

	"hexagonal-demo/modules/playlist"
	playlistadapters "hexagonal-demo/modules/playlist/adapters"
	"hexagonal-demo/modules/song"
	songadapters "hexagonal-demo/modules/song/adapters"
	"hexagonal-demo/modules/user"
	useradapters "hexagonal-demo/modules/user/adapters"
)

// RegisterModules wires all modules in two phases:
//  1. create all module references — pointers are stable from here on, so
//     adapters can capture them regardless of dependency cycles
//  2. initialize each module and mount its routes — order does not matter,
//     nothing dereferences a module until the server starts serving
func RegisterModules(db *sql.DB, mux *http.ServeMux) {
	songModule := &song.Module{}
	playlistModule := &playlist.Module{}
	userModule := &user.Module{}

	registerSongModule(db, mux, songModule)
	registerPlaylistModule(db, mux, playlistModule, songModule)
	registerUserModule(db, mux, userModule, playlistModule)
}

func registerSongModule(db *sql.DB, mux *http.ServeMux, module *song.Module) {
	module.Initialize(songadapters.NewSqlAdapter(db))
	songadapters.NewHTTPAdapter(module.Facade()).RegisterRoutes(mux)
}

func registerPlaylistModule(db *sql.DB, mux *http.ServeMux, module *playlist.Module, songs *song.Module) {
	module.Initialize(
		playlistadapters.NewSqlAdapter(db),
		playlistadapters.NewSongModuleAdapter(songs.Facade()),
	)
	playlistadapters.NewHTTPAdapter(module).RegisterRoutes(mux)
}

func registerUserModule(db *sql.DB, mux *http.ServeMux, module *user.Module, playlists *playlist.Module) {
	module.Initialize(
		useradapters.NewSqlAdapter(db),
		useradapters.NewPlaylistModuleAdapter(playlists),
	)
	useradapters.NewHTTPAdapter(module).RegisterRoutes(mux)
}
