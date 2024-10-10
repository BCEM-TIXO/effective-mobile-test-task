package migration

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) error {
	q := `CREATE TABLE IF NOT EXISTS song (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link VARCHAR(512),
		is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	_, err := pool.Exec(context.Background(), q)
	return err
}
