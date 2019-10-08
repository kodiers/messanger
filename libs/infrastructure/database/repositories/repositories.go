package repositories

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	SelectQuery()
	CreateInsertQuery()
}

func SelectQuery(db *sql.DB, query string) (*sql.Rows, error) {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error on executing query: ", query, err)
		return nil, err
	}
	return rows, nil
}

func RunCreateInsertQuery(db *sql.DB, query string) error {
	_, err := db.Query(query)
	if err != nil {
		fmt.Println("Error on executing query: ", query, err)
		return err
	}
	return nil
}
