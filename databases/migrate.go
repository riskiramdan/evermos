package databases

import (
	"database/sql"
	"fmt"
	"log"

	rice "github.com/GeertJohan/go.rice"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/riskiramdan/evermos/config"
)

// MigrateUp migrates the database up
func MigrateUp() {
	// Setup the database
	//
	cfg, err := config.GetConfiguration()
	if err != nil {
		log.Fatal("error when getting configuration: ", err)
	}

	db, err := sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		log.Fatal("error when open postgres connection: ", err)
	}

	// Setup the source driver
	//
	sourceDriver := &RiceBoxSource{}
	sourceDriver.PopulateMigrations(rice.MustFindBox("./migrations"))
	if err != nil {
		log.Fatal("error when creating source driver: ", err)
	}

	// Setup the database driver
	//
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("error when creating postgres instance: ", err)
	}

	m, err := migrate.NewWithInstance(
		"go.rice", sourceDriver,
		"postgres", driver)

	if err != nil {
		log.Fatal("error when creating database instance: ", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			log.Fatal("error when migrate up: ", err)
		}
	}
	fmt.Println("success migrate databases")

	defer m.Close()
}
