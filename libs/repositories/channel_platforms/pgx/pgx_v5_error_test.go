package pgx

import pgxv5pgconn "github.com/jackc/pgx/v5/pgconn"

func newPgxV5UniqueViolationError() error {
	return &pgxv5pgconn.PgError{Code: "23505"}
}
