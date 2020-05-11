package main

import (
	"encoding/json"
	"net/http"
)

func addRating(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden! ", http.StatusForbidden)
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
			UPDATE comments SET rating=(?) WHERE mvid=(?) AND uid=(?)`
			_, err := db.Exec(sqlStatement, com.Rating, com.MVID, com.UID)
			if err != nil {
				panic(err)
			}
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

}
