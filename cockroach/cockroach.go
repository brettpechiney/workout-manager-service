package cockroach

import (
	"database/sql"

	_ "github.com/lib/pq" // Import underlying database driver.
	"github.com/pkg/errors"
)

// Cockroach provides application-level context to the database handle.
type Cockroach struct {
	db *sql.DB
}

// NewCockroach creates a database object, associates it with the Postgres
// driver.
func NewCockroach(dataSource string) (Cockroach, error) {
	conn, err := sql.Open("postgres", dataSource)
	if err != nil {
		return Cockroach{}, errors.Wrap(err, "failed to open database connection")
	}
	return Cockroach{conn}, nil
}

// Close closes the underlying sql.DB
func (m Cockroach) Close() error {
	return m.db.Close()
}
