package main

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tbhsuseumr1"
	dbname   = "skool"
)

var db *sql.DB

func connect() {
	var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func InsertUser(name string, userpassword string) {
	connect()
	sqlStatement := `
INSERT INTO users (username,password)
VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, name, userpassword)
	if err != nil {
		panic(err)
	}
}

func isUserCreds(name string, userpassword string) (stt bool) {
	connect()
	sqlstt := "select * from users where username=$1 and password=$2;"
	var tmpname, tmppassword string
	var index int
	row := db.QueryRow(sqlstt, name, userpassword)
	switch err := row.Scan(&index, &tmpname, &tmppassword); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		panic(err)
	}

}

func setSession(id string, name string) {
	connect()
	sqlStatement := `
INSERT INTO sessions (username,uuid)
VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, name, id)
	if err != nil {
		panic(err)
	}
}

func getSession(id string) (name string) {
	connect()
	sqlstt := "select username from sessions where uuid=$1;"
	var tmpname string
	row := db.QueryRow(sqlstt, id)
	switch err := row.Scan(&tmpname); err {
	case sql.ErrNoRows:
		return "no such"
	case nil:
		return tmpname
	default:
		panic(err)
	}

}

func isUser(name string) bool {
	connect()
	sqlstt := "select username from users where username=$1;"
	var tmpname string
	row := db.QueryRow(sqlstt, name)
	switch err := row.Scan(&tmpname); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		panic(err)
	}

}

func setPost(name string, post string, post_time time.Time) {
	connect()
	sqlStatement := `
INSERT INTO posts (username,post_time, post)
VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, name, post_time, post)
	if err != nil {
		panic(err)
	}

}

func getPosts() []string {
	connect()
	//sqlstt := `select * from posts`
	//under construction
}
