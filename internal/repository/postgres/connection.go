package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/sementrof/prod1/internal/models"
)

type Config struct{ Host, Port, User, Password, Dbname string }

func Connection(conf Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://"+conf.User+":"+conf.Password+"@"+conf.Host+":"+conf.Port+"/"+conf.Dbname)
	if err != nil {
		return nil, models.ErrConnectionDb
	}
	return conn, nil
}
