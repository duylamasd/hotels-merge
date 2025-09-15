package config

import (
	sqlc "github.com/duylamasd/hotels-merge/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStore struct {
	Queries  sqlc.Querier
	ConnPool *pgxpool.Pool
}

func NewDBStore(connPool *pgxpool.Pool) *DBStore {
	return &DBStore{
		Queries:  sqlc.New(connPool),
		ConnPool: connPool,
	}
}
