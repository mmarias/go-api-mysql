package main

import (
	"fmt"
	"log"
	"net/http"
	//"io/ioutil"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_"github.com/go-sql-driver/mysql"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users interface{}
	db, err := sql.Open("mysql", "root:12345678@tcp(localhost:3306)/go")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db is connected")
	}
	defer db.Close()

	resp, err := db.Prepare("SELECT * FROM user")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Close()

	json.NewEncoder(w).Encode(resp)
	fmt.Println(resp)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", getUsers).Methods("GET")
	router.HandleFunc("/user", addUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func CreateCon() *sql.DB {
	db, err := sql.Open("mysql", "root:12345678@tcp(localhost:3306)/go")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db is connected")
	}
	defer db.Close()

	err = db.Ping()
	fmt.Println(err)
	if err != nil {
		fmt.Println("MySQL db is not connected")
		fmt.Println(err.Error())
	}
	return db
}