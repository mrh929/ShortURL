package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

func dbInit() (err error) {
	fmt.Println("dbInit...")

	sqlSrvURL := fmt.Sprintf("root:%s@tcp(%s:%s)/", SQL_ROOT_PASSWD, SQL_HOST, SQL_PORT)

	db, err = sql.Open("mysql", sqlSrvURL)
	if err != nil {
		return
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", SQL_DATABASE_NAME))
	if err != nil {
		return
	}

	db.Close()

	bytes, err := ioutil.ReadFile("./sql/table.sql")
	if err != nil {
		return
	}

	db, err = sql.Open("mysql", sqlSrvURL+SQL_DATABASE_NAME)
	if err != nil {
		return
	}

	_, err = db.Exec(string(bytes))
	if err != nil {
		return
	}
	fmt.Println("dbInit finished")
	return
}

func urlSelect(s_key string) (r_url string, err error) {
	err = db.Ping()
	if err != nil {
		return
	}
	rows, err := db.Query("SELECT r_url FROM urltable WHERE s_key = ?", s_key)

	for rows.Next() {
		err = rows.Scan(&r_url)
		if err != nil {
			return
		}
		log.Println(s_key, r_url)
		return
	}
	err = rows.Err()
	r_url = ""
	rows.Close()
	return
}

func urlInsert(s_key string, r_url string) (err error) {
	err = db.Ping()
	if err != nil {
		return
	}
	rows, err := db.Query("INSERT INTO urltable VALUES(?, ?, null)", s_key, r_url)
	rows.Close()
	return
}
