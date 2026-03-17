package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Database struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	DatabaseName string `json:"databaseName"`
	Maxconns     int    `json:"maxConns"`
	Minconns     int    `json:"minConns"`
	once         sync.Once
}

func (dbCreds *Database) InitDb(ctx context.Context, log zerolog.Logger) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	dbCreds.once.Do(func() {
		databaseUrl := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbCreds.Username, dbCreds.Password, dbCreds.DatabaseName, dbCreds.Host, dbCreds.Port)

		dbConfig, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			err = fmt.Errorf("Error while initiating the database pool : %w", err)
			return
		}

		// setting the minimum and maximum connections
		dbConfig.MinConns = int32(dbCreds.Minconns)
		dbConfig.MaxConns = int32(dbCreds.Maxconns)

		// creating the pool with the config
		pool, err = pgxpool.NewWithConfig(ctx, dbConfig)
		if err != nil {
			log.Error().Err(err).Msg("Error initiating the database pool : ")
			err = fmt.Errorf("Error initiating the database pool : %w", err)
			return
		}
	})

	// ping the database to check the connection
	if err := pool.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Database connection ping failed")
		return nil, fmt.Errorf("Database connection ping failed : %w", err)
	}

	return pool, err
}

func CreateTables(ctx context.Context, db *pgxpool.Pool, log zerolog.Logger, tables map[string]string) {
	for tableName, createQuery := range tables {
		_, err := db.Exec(ctx, createQuery)

		if err != nil {
			log.Info().Err(err).Msg("Error creating table : " + tableName)
		} else {
			log.Info().Msg("Table " + tableName + " created/exists")
		}
	}
}
