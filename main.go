package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	_"github.com/go-sql-driver/mysql"
)

type ApiResponse struct {
	Result string `json:"result"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Status string `json:"status"`
	Sex string `json:"sex"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User

	db, err := sql.Open("mysql", "root:12345678@tcp(localhost:3306)/go")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db is connected")
	}

	defer db.Close()

	resp, err := db.Query("SELECT username, password, status, sex FROM user")

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Close()
	
	for resp.Next() {
		var user User

		err = resp.Scan(&user.Username, &user.Password, &user.Status, &user.Sex)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
		}

		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:12345678@tcp(localhost:3306)/go")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{"Error on connecto to DB"})
	}

	defer db.Close()

	reqBody, err := ioutil.ReadAll(r.Body)
	

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{"Error on request parameters"})
	}

	var newUser User

	json.Unmarshal(reqBody, &newUser)
	
	if newUser.Username == "" || newUser.Password == "" || newUser.Status == "" || newUser.Sex == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{"Error on request parameters"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{"Error on hash password"})
	}

	insert, err := db.Query(""+
	"INSERT INTO user (username, password, status, sex)"+
	"VALUES ('"+newUser.Username+"','"+string(hash)+"','"+newUser.Status+"','"+newUser.Sex+"')")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{"Error on create user"})
	}

	defer insert.Close()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ApiResponse{"User created"})
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