package main

import (
	"log"

	"github.com/riskiramdan/evermos/databases"
	"github.com/riskiramdan/evermos/seeder"
)

func main() {
	databases.MigrateUp()
	err := seeder.SeedUp()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
}
