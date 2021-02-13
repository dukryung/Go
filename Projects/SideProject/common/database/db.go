package database

import (
	"database/sql"
	"log"
)

type DBHandler interface {
	ReadProjectList()
}

type mariadbHandler struct {
	db *sql.DB
}

type ProjectDB struct {
}

type DB interface {
	NewMariaDBHandler() *mariadbHandler
}

func MakeDBHandler(database DB) DBHandler {
	return database.NewMariaDBHandler()
}

func (p *ProjectDB) NewMariaDBHandler() *mariadbHandler {
	database, err := sql.Open("mysql", "delis:@shud")

	if err != nil {
		log.Println("[LOG] database open err : ", err)
		return nil
	}

	return &mariadbHandler{db: database}
}

func (m *mariadbHandler) ReadProjectList() {

}

func (m *mariadbHandler) CreateProjectListTable() {

}
