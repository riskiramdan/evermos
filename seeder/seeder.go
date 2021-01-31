package seeder

import (
	"database/sql"
	"fmt"

	"github.com/riskiramdan/evermos/config"
)

// SeedUp seeding the database
func SeedUp() error {
	cfg, err := config.GetConfiguration()
	if err != nil {
		return fmt.Errorf("error when getting configuration: %s", err)
	}

	db, err := sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		return fmt.Errorf("error when open postgres connection: %s", err)
	}
	defer db.Close()

	_, err = db.Exec(`
	insert into "user" ("name", "email", "password", "createdAt", "updatedAt") values
	('Admin evermos', 'admin', '$2a$10$CCWkaZ1UedUQeGtABKwEqepsNTBR1Rp.b4UlFCuvtkbGmfX3rk4SC', now(), now()),
	('author', 'author@evermos.com', '$2a$10$CCWkaZ1UedUQeGtABKwEqepsNTBR1Rp.b4UlFCuvtkbGmfX3rk4SC', now(), now());
	`)
	if err != nil {
		return err
	}

	return nil
}
