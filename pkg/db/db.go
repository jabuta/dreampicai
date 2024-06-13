package db

import (
	"context"
	"embed"
	"os"

	"github.com/jabuta/dreampicai/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type DBConfig struct {
	DB            database.Queries
	PgxConneciton pgx.Conn
	Ctx           context.Context
}

var Conf = &DBConfig{}

func InitDatabase(migrations embed.FS) error {
	var err error
	Conf.Ctx = context.Background()
	pgxConneciton, err := pgx.Connect(context.Background(), os.Getenv("DB_STRING"))
	if err != nil {
		return err
	}
	Conf.PgxConneciton = *pgxConneciton

	// Run migrations
	db := stdlib.OpenDB(*Conf.PgxConneciton.Config())

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	// if err := goose.Reset(db, "sql/schema"); err != nil {
	// 	return err
	// }
	if err := goose.Up(db, "sql/schema"); err != nil {
		return err
	}

	Conf.DB = *database.New(&Conf.PgxConneciton)
	return nil
}
