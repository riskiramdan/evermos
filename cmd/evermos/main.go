package main

import (
	"fmt"
	"log"
	testUser "os/user"

	"github.com/jmoiron/sqlx"
	"github.com/riskiramdan/evermos/config"
	"github.com/riskiramdan/evermos/databases"
	"github.com/riskiramdan/evermos/internal/data"
	internalhttp "github.com/riskiramdan/evermos/internal/http"
	"github.com/riskiramdan/evermos/internal/product"
	productPg "github.com/riskiramdan/evermos/internal/product/postgres"
	"github.com/riskiramdan/evermos/internal/user"
	userPg "github.com/riskiramdan/evermos/internal/user/postgres"
)

// InternalServices represents all the internal domain services
// that will be used by payfazz-commerce httpserver
type InternalServices struct {
	userService    user.ServiceInterface
	productService product.ServiceInterface
}

func buildInternalServices(db *sqlx.DB, config *config.Config) *InternalServices {
	userPostgresStorage := userPg.NewPostgresStorage(
		data.NewPostgresStorage(db, "user", user.User{}),
	)
	userService := user.NewService(userPostgresStorage)
	productPostgresStorage := productPg.NewPostgresStorage(
		data.NewPostgresStorage(db, "product", product.Product{}),
	)
	productService := product.NewService(productPostgresStorage)
	return &InternalServices{
		userService:    userService,
		productService: productService,
	}
}

func main() {

	usr, err := testUser.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usr.Username)

	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}
	db, err := sqlx.Open("postgres", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	defer db.Close()
	dataManager := data.NewManager(db)
	internalServices := buildInternalServices(db, config)
	// Migrate the db
	databases.MigrateUp()

	s := internalhttp.NewServer(
		internalServices.userService,
		internalServices.productService,
		dataManager,
		config,
	)
	s.Serve()
}
