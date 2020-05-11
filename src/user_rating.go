package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UserRating struct {
	Name     string
	Rating   int
	Comments string
	Count    int `json:"-"`
	UID      int `json:"-"`
}

func getUserRating(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden! No user logged in", http.StatusForbidden)
		return
	}
	user1 := session.Values["uid"]
	getmov := GetMovie{}
	rows, err := db.Query("SELECT m.mvname,c.uid,c.mvid,c.rating,c.comments FROM movie m JOIN comments c ON (m.mvid=c.mvid)")
	CheckErr(err)
	for rows.Next() {
		err = rows.Scan(&getmov.Name, &getmov.UID, &getmov.MVID, &getmov.Rating, &getmov.Comments)
		CheckErr(err)
		if getmov.UID == user1 {
			out, err := json.Marshal(getmov)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			fmt.Fprintf(w, string(out))
		}
	}

}
