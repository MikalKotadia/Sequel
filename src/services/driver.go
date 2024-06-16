package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sequel/main/utils"

	"github.com/georgysavva/scany/dbscan"
	_ "github.com/georgysavva/scany/dbscan"
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
	if !isSupported {
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

func (con Connection) GetTables() ([]string, error) {
	var query_string string
	var table_name string

	println(con.variation)
	switch con.variation {
	case "postgres":
		table_name = "tablename"
		query_string = fmt.Sprintf("SELECT %s FROM pg_catalog.pg_tables WHERE schemaname='public';", table_name)
	default:
		return nil, errors.ErrUnsupported
	}

	result, err := con.executeQuery(query_string)
	if err != nil {
		return nil, err
	}
    reducer := func (item map[string]string) string {
        return item[table_name]
    }

    formatted_result := utils.Map(result, reducer)
	return formatted_result, nil

}

func (con Connection) DropTable(table_name string) error {
    query_string := fmt.Sprintf("DROP TABLE %s;", table_name)
    if _, err := con.executeQuery(query_string); err != nil {
        return err
    }

    return nil
}

func (con Connection) executeQuery(query string) ([]map[string]string, error) {
	rows, err := con.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

    var result []map[string]string // TODO: need to check if this will cast
	if err := dbscan.ScanAll(&result, rows); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
