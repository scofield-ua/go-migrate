package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/scofield-ua/go-migrate/config"
)

func ConnectPostgreSQL(config *config.Config) (*pgx.Conn, error) {
	conn, err := ConnectToDatabase(config.DB.Host, config.DB.Username, config.DB.Password, config.DB.Database)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
