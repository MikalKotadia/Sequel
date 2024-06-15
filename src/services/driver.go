package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sequel/main/utils"

	_ "github.com/lib/pq"
)

type Connection struct {
	variation string
	host      string
	port      int
	user      string
	password  string
	dbname    string
	db        *sql.DB
}

type DbDriver interface {
	createConnection(string, int, string, string, string) error
	getTables() []string
}

func (con *Connection) CreateConnection(variation string, host string, port int, user string, password string, dbname string) error {
	_, isSupported := utils.GetSupportedVariants()[variation]
	if (!isSupported) {
		return errors.New("Database of type " + variation + " is not supported")
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open a connection to the database
	db, err := sql.Open(variation, connStr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Attempt to ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Successfully connected to the database")
	con.variation = variation
	con.host = host
	con.port = port
	con.user = user
	con.password = password
	con.dbname = dbname
	con.db = db

    return nil
}

func (con Connection) GetTables() (*sql.Rows, error) {
	var query_string string

    println(con.variation)
	switch con.variation {
	case "postgres":
		query_string = "SELECT 'tablename' FROM pg_catalog.pg_tables WHERE schemaname='public';"
	default:
		return nil, errors.ErrUnsupported
	}

    tables, err := con.db.Query(query_string)
    if (err != nil) {
        return nil, err
    }

    // we are making it here, but the result is nil....
	return tables, nil

}
