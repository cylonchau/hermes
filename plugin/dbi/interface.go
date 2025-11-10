package dbi

import "database/sql"

type DatabaseInterface interface {
	Connect() (*sql.DB, error)
	GetConnectionString() string
	GetDriverName() string
	EscapeIdentifier(identifier string) string
}
