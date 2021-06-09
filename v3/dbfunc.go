package main

import (
	"database/sql"
	"fmt"
)

func InsertUser(name string, userpassword string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	sqlStatement := `
INSERT INTO users (username,password)
VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, name, userpassword)
	if err != nil {
		panic(err)
	}
}
