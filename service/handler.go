package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	db    *sqlx.DB
	store sessions.Store
}

func NewHandlers(db *sqlx.DB, store sessions.Store) *Handlers {
	return &Handlers{
		db:    db,
		store: store,
	}
}
