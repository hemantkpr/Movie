package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Comment struct {
	UID      int
	MVID     int
	Rating   float32
	Comments string
}
type CommentCheck struct {
	UID      int
	MVID     int
	Rating   float32
	Comments string
}
type Moviee struct {
	MVID   int
	MVNAME string
}

func addComment(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden! No user logged in", http.StatusForbidden)
		return
	}

	com := Comment{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&com)

	defer handlepanic()

	mv := CommentCheck{}
	rows, err := db.Query("SELECT * FROM comments")
	CheckErr(err)
	for rows.Next() {
		err = rows.Scan(&mv.UID, &mv.MVID, &mv.Rating, &mv.Comments)
		CheckErr(err)

		if com.MVID == mv.MVID && com.UID == mv.UID {

			sqlStatement := `
			UPDATE comments SET comments=(?) WHERE mvid=(?) AND uid=(?)`
			_, err := db.Exec(sqlStatement, com.Comments, com.MVID, com.UID)
			if err != nil {
				panic(err)
			}
			fmt.Println("Updated")
			return
		}

	}

	sqlStatement := `
			INSERT INTO comments (uid, mvid, rating,comments)
			VALUES ((?), (?), (?), (?))`
	_, errr := db.Exec(sqlStatement, com.UID, com.MVID, com.Rating, com.Comments)
	if errr != nil {
		panic(errr)
	}
	fmt.Println("Updated")

}

func handlepanic() {
	if a := recover(); a != nil {
		fmt.Println("Movie not available!")
	}
}
