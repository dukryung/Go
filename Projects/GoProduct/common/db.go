package common

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DBfilepath struct {
	Productdbfilepath string
	Userdbfilepath    string
}

type DB interface {
	NewSqliteHandler() *sqliteHandler
}

type ProductDB struct {
	filepath string
}
type UserDB struct {
	filepath string
}

type Products struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SessionID int    `json:"sessionid"`
}

type Users struct {
	ID        int    `json:"id"`
	Name      string `json:"name:`
	SessionID int    `json:"sessionid"`
}

type DBHandler interface {
	createProducts(string, string)
	readProducts(string) []*Products
	updateProducts(string, string)
	deleteProducts(string, string)
	createUsers(*userinfo)
	readUsers(int, string) []*Users
	updateUsers(int, string)
	deleteUsers(string, int)

	Close()
}

type sqliteHandler struct {
	db *sql.DB
}

func NewDBHandler(database DB) DBHandler {
	return database.NewSqliteHandler()
}

func CreateDBFolder() {
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		err := os.MkdirAll("./db", 0777)
		if err != nil {
			fmt.Println("[LOG] Create folder Error :", err)
			return
		}
	}
}

func (pdb *ProductDB) NewSqliteHandler() *sqliteHandler {
	CreateDBFolder()

	database, err := sql.Open("sqlite3", pdb.filepath)
	if err != nil {
		panic(err)
	}
	stmt, err := database.Prepare(`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sessionid STRING,
			title TEXT,
			desc TEXT,
			createdat DATETIME
			);
			CREATE INDEX IF NOT EXISTS sessionidIndexOnproducts ON products (
				sessionid ASC
			);
			`)

	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	stmt.Exec()

	return &sqliteHandler{db: database}
}

func (udb *UserDB) NewSqliteHandler() *sqliteHandler {
	database, err := sql.Open("sqlite3", udb.filepath)
	if err != nil {
		panic(err)
	}

	stmt, err := database.Prepare(`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			sessionid STRING,
			userid STRING,
			password STRING,
			email STRING
			);
			CREATE INDEX IF NOT EXISTS sessionidIndexOnusers ON users (
				sessionid ASC
			);
			`)

	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	stmt.Exec()
	return &sqliteHandler{db: database}
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}
