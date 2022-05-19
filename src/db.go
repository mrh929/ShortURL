package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

func dbInit() (err error) {
	log.Info("dbInit...")

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
	log.Info("dbInit finished")
	return
}

func urlSelect(s_key string) (r_url string, err error) {
	rows, err := db.Query("SELECT r_url FROM urltable WHERE s_key = ?", s_key)
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&r_url)
		if err != nil {
			return
		}
	}
	err = rows.Err()
	rows.Close()
	return
}

func urlInsert(s_key string, r_url string) (err error) {
	rows, err := db.Query("INSERT INTO urltable VALUES(?, ?, null)", s_key, r_url)
	if err != nil {
		return
	}
	rows.Close()
	return
}
