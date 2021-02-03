package apiserver

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/UrcaDeLima/backend_golang_journal/internal/app/store/sqlstore"
	_ "github.com/lib/pq" // ...
)

// Start ...
func Start() error {
	db, err := newDB(os.Getenv("DB_URL"))
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	srv := newServer(store)

	return http.ListenAndServe(os.Getenv("DB_PORT"), srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open(os.Getenv("DB_DIALECT"), dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
