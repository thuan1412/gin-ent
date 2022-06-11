package helpers

import (
	"entgo.io/ent/dialect"
	"gin-ent/ent"
)

// GetDb return an ent client
func GetDb() (*ent.Client, error) {
	return ent.Open(dialect.Postgres,
		"host=localhost port=5432 user=postgres dbname=gin-ent password=1 sslmode=disable")
}
