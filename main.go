package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// load .env if present; real environment variables take precedence
	_ = godotenv.Load()

	config := ConfigFromEnv()

	db, err := openDB(config.DB)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer func() { _ = db.Close() }()
	log.Printf("connected to mysql at %s", config.DB.HOST)

	mux := http.NewServeMux()

	RegisterModules(db, mux)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "db unreachable", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":" + config.HTTP.PORT,
		Handler: mux,
	}
	log.Printf("listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("http: %v", err)
	}
}

func openDB(cfg DBConfig) (*sql.DB, error) {
	mc := mysql.NewConfig()
	mc.Net = "tcp"
	mc.Addr = net.JoinHostPort(cfg.HOST, cfg.PORT)
	mc.DBName = cfg.DATABASE
	mc.User = cfg.USERNAME
	mc.Passwd = cfg.PASSWORD
	mc.ParseTime = true

	db, err := sql.Open("mysql", mc.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
