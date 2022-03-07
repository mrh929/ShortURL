package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

func dbInit() (err error) {
	fmt.Println("dbInit...")

	db, err = sql.Open("mysql", "root:test@tcp(127.0.0.1:3306)/")
	if err != nil {
		return
	}
	db.Exec("CREATE DATABASE IF NOT EXISTS db;")
	db.Close()

	bytes, err := ioutil.ReadFile("./sql/table.sql")
	if err != nil {
		return
	}

	db, err = sql.Open("mysql", "root:test@tcp(127.0.0.1:3306)/db")
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
	defer rows.Close()

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
	return
}

func urlInsert(s_key string, r_url string) (err error) {
	err = db.Ping()
	if err != nil {
		return
	}
	rows, err := db.Query("INSERT INTO urltable VALUES(?, ?, null)", s_key, r_url)
	defer rows.Close()
	return
}
