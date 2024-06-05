package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/agnate/qlikrestapi/api/router"
	"github.com/agnate/qlikrestapi/config"
	"github.com/agnate/qlikrestapi/internal/migrator"
)

// "driver://user:pass@host:port/dbName?sslmode=[enable|disable]"
const dbConnection = "%s://%s:%s@%s:%d/%s?sslmode=%s"

func main() {
	// Load environment config.
	c := config.New()

	// Connect to database.
	port, _ := strconv.Atoi(c.Database.Port)
	connStr := fmt.Sprintf(dbConnection, c.Database.Driver, c.Database.Username, c.Database.Password, c.Database.Host, port, c.Database.DatabaseName, c.Database.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Create migrator and run it.
	migrator := migrator.New("./migrations/", c.Database.DatabaseName)
	err = migrator.Run(db)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Add auth middleware between http and router.

	// Initialize API router.
	router := router.New(db)

	// Serve API router.
	apiPort, _ := strconv.Atoi(c.API.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", apiPort), router.NewHandler()))
}
