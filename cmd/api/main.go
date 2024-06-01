package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/agnate/qlikrestapi/api/entity/message"
	"github.com/agnate/qlikrestapi/config"
)

// "driver://user:pass@host:port/dbName?sslmode=[enable|disable]"
const dbConnection = "%s://%s:%s@%s:%d/%s?sslmode=%s"

func main() {
	// Load environment config
	c := config.New()

	// Connect to database
	port, _ := strconv.Atoi(c.Database.Port)
	connStr := fmt.Sprintf(dbConnection, c.Database.Driver, c.Database.Username, c.Database.Password, c.Database.Host, port, c.Database.DatabaseName, c.Database.SSLMode)
	log.Print(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM messages")
		if err != nil {
			log.Fatalln(err)
		}
		var msg message.Message
		var found bool = false
		for rows.Next() {
			rows.Scan(&msg.Username, &msg.CreateDate, &msg.Message, &msg.Palindrome)
			//log.Print(order)
			fmt.Fprintf(w, msg.Username+" - "+msg.CreateDate+" - Palindrome? "+msg.Palindrome+" - Message: '"+msg.Message+"'\n")
			found = true
		}
		if !found {
			fmt.Fprintf(w, "No messages found.")
		}
	})
	apiPort, _ := strconv.Atoi(c.API.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", apiPort), nil))
}
