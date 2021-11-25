package db

// schema.go provides data models in DB
import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt null.Time `db:"deleted_at"`
}

// Task corresponds to a row in `tasks` table
type Task struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	Title     string    `db:"title"`
	IsDone    bool      `db:"is_done"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt null.Time `db:"deleted_at"`
}
