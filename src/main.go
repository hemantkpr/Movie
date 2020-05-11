package main

import (
	"database/sql"
	"fmt"
	_ "fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	db    *sql.DB
	keyy  = []byte("super-secret-key")
	store = sessions.NewCookieStore(keyy)
)

type User struct {
	UID      int
	Username string
	Password int
	Role     string
}

func init() {
	//DB Connection
	var err error
	fmt.Println("Getting started")
	connString := "root:hemant7830@tcp(127.0.0.1:3306)/movie"
	db, err = sql.Open("mysql", connString)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected")

}

func main() {
	//API call
	r := mux.NewRouter()
	r.HandleFunc("/login/{user}", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/movie/{movie}", getMovieDetails)
	r.HandleFunc("/userating", getUserRating)
	r.HandleFunc("/comment", addComment)
	r.HandleFunc("/rating", addRating)

	log.Fatal(http.ListenAndServe(":8080", r))
	defer db.Close()
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	usr := User{}
	vars := mux.Vars(r)
	key := vars["user"]
	rows, err := db.Query("SELECT uid,uname,role FROM user")
	CheckErr(err)
	for rows.Next() {
		err = rows.Scan(&usr.UID, &usr.Username, &usr.Role)
		CheckErr(err)
		if usr.Username == key && usr.Role == "admin" {
			session.Values["authenticatedAdmin"] = true
			session.Values["uid"] = usr.UID
			session.Values["user"] = usr.Username
			session.Values["role"] = usr.Role
			session.Save(r, w)
			fmt.Println("admin login")
			return
		} else if usr.Username == key && usr.Role == "user" {
			session.Values["authenticated"] = true
			session.Values["uid"] = usr.UID
			session.Values["user"] = usr.Username
			session.Values["role"] = usr.Role
			session.Save(r, w)
			fmt.Println("user login")
			return
		} else {
			session.Values["authenticated"] = false
			session.Save(r, w)
			fmt.Println("User logout")
		}
	}

}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticatedAdmin"] = false
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
