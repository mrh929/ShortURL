package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var db *sql.DB
var log = logrus.New()
var SRV_PASSWD string
var SRV_HOST string
var SRV_PORT string
var SRV_PROTO string
var SRV_BASE_PATH string
var SQL_ROOT_PASSWD string
var SQL_HOST string
var SQL_PORT string
var SQL_DATABASE_NAME string

func main() {
	allInit()

	log.Info("Serving Short URLs...")
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler) // add handlers
	router.HandleFunc("/shorten", shortenHandler)
	router.HandleFunc("/{key}", urlHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", SRV_HOST, SRV_PORT), router)) // start listening
}

func allInit() {
	envSet() // set env

	err := dbInit() // init db
	if err != nil {
		log.Fatal(err.Error())
	}
}

func envSet() {
	SRV_PASSWD = os.Getenv("SRV_PASSWD") // init server password
	if SRV_PASSWD == "" {
		SRV_PASSWD = "123456!"
	}
	log.Info("SRV_PASSWD = ", SRV_PASSWD)

	SRV_HOST = os.Getenv("SRV_HOST") // init server host
	if SRV_HOST == "" {
		SRV_HOST = "0.0.0.0"
	}
	log.Info("SRV_HOST = ", SRV_HOST)

	SRV_PORT = os.Getenv("SRV_PORT") // init server port
	if SRV_PORT == "" {
		SRV_PORT = "8000"
	}
	log.Info("SRV_PORT = ", SRV_PORT)

	SRV_PROTO = os.Getenv("SRV_PROTO") // init client access protocal (e.g. http/https)
	if SRV_PROTO == "" {
		SRV_PROTO = "http"
	} else if SRV_PROTO != "https" && SRV_PROTO != "http" {
		log.Fatal("server protocol not supported")
	}
	log.Info("SRV_PROTO = ", SRV_PROTO)

	SRV_BASE_PATH = os.Getenv("SRV_BASE_PATH") // init url base path (e.g. exam.ple/foo/bar/ABCDE)
	SRV_BASE_PATH = strings.TrimRight(SRV_BASE_PATH, "/")
	log.Info("SRV_BASE_PATH = ", SRV_BASE_PATH)

	SQL_ROOT_PASSWD = os.Getenv("SQL_ROOT_PASSWD") // init sql root password
	if SQL_ROOT_PASSWD == "" {
		SQL_ROOT_PASSWD = "test"
	}
	log.Info("SQL_ROOT_PASSWD = ", SQL_ROOT_PASSWD)

	SQL_HOST = os.Getenv("SQL_HOST") // init sql server host
	if SQL_HOST == "" {
		SQL_HOST = "127.0.0.1"
	}
	log.Info("SQL_HOST = ", SQL_HOST)

	SQL_PORT = os.Getenv("SQL_PORT") // init sql server port
	if SQL_PORT == "" {
		SQL_PORT = "3306"
	}
	log.Info("SQL_PORT = ", SQL_PORT)

	SQL_DATABASE_NAME = os.Getenv("SQL_DATABASE_NAME") // init sql database name
	if SQL_DATABASE_NAME == "" {
		SQL_DATABASE_NAME = "db"
	}
	log.Info("SQL_DATABASE_NAME = ", SQL_DATABASE_NAME)
}
