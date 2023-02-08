package dbrepo

import (
	"database/sql"
	"github.com/tnthanh47/Booking/internal/config"
	"github.com/tnthanh47/Booking/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(a *config.AppConfig, conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
