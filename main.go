package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var db *sql.DB
var log = logrus.New()
var PASSWD string

func initAll() {
	err := dbInit() // init db
	if err != nil {
		log.Fatal(err.Error())
	}

	pass, err := ioutil.ReadFile("./PASSWD") // init password
	PASSWD = string(pass)
}

func main() {
	initAll()

	fmt.Println("Serving Short URLs...")
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler) // add handlers
	router.HandleFunc("/shorten", shortenHandler)
	router.HandleFunc("/{key}", urlHandler)

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", router)) // start listening
}
