package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

/*

This file contains an interface that unions pgxpool.Pool, pgx.Conn, pgxpool.Tx

Because in some cases the source of data does not matter, but might be from all of those.

*/

type Connection interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func GetConnectionFromContextOrDefault(ctx context.Context, defaultConn Connection) Connection {
	c := ctx.Value("connection").(Connection)
	if c != nil {
		return c
	}
	return defaultConn
}
