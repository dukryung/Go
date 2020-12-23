package common

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Products struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SessionID int    `json:"sessionid"`
}

type DBHandler interface {
	getProducts(string) []*Products
	addProducts(string, string)
	Close()
}

type sqliteHandler struct {
	db *sql.DB
}

func NewDBHandler(filepath string) DBHandler {
	return NewSqliteHandler(filepath)
}

func NewSqliteHandler(filepath string) *sqliteHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	stmt, err := database.Prepare(`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			sessionid STRING,
			createdat DATETIME,
			);
			CREATE INDEX IF NOT EXISTS sessioniddIndexOnProducts ON Products (
				sessionid ASC
			)
			`)

	if err != nil {
		panic(err)
	}
	stmt.Exec()

	return &sqliteHandler{db: database}
}
func (s *sqliteHandler) Close() {
	s.db.Close()
}

func (s *sqliteHandler) getProducts(sessionid string) []*Products {
	// read html file
	productlist := []*Products{}
	rows, err := s.db.Query("SELECT id, name, sessionid FROM products WHERE sessionid=?", sessionid)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var product Products
		rows.Scan(&product.ID, &product.Name, &product.SessionID)
		productlist = append(productlist, &product)
	}
	return productlist
}

func (s *sqliteHandler) addProducts(name string, sessionid string) {
	stmt, err := s.db.Prepare("INSERT INTO products (name, sessionid, createdat)  VALUES(??, datetime('now')")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, sessionid)
	if err != nil {
		panic(err)
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		fmt.Println("[ERR] RowsAffected err : ", err)
		panic(err)
	}

	if rowcnt != 1 {
		fmt.Println("[LOG] anything is not affected")
	}

}
