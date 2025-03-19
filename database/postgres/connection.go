package postgres

import (
	"fmt"
	"labyrinth/config"
)

func getConnectionStr(conf *config.Config) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s",
		conf.PostgreSQL.Username, conf.PostgreSQL.Password,
		conf.PostgreSQL.DatabaseName, conf.PostgreSQL.SSLMode,
		conf.PostgreSQL.Host,
	)
}

func NewPostgres(conf *config.Config) *Postgres {
	newPostgres := &Postgres{conn: getConnectionStr(conf)}
	return newPostgres
}
