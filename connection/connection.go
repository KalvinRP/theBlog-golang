package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn
var err error

func DatabaseConnect() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	databaseUrl := "postgres://postgres:3901@localhost:5432/theBlog-db"

	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Println("Connection failed:", err)
		os.Exit(1)
	}

	fmt.Println("Database connected")
}
