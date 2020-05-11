package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type GetMovie struct {
	Name     string
	UID      int
	MVID     int
	Rating   float32
	Comments string
}
type Movie struct {
	Name     string
	Count    int
	AVG      float32
	Comments string
}

func getMovieDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["movie"]

	getmov := Movie{}

	sqlquery, err := db.Query("SELECT m.mvname,AVG(c.rating), COUNT(c.rating),c.comments FROM comments c JOIN movie m ON (m.mvid=c.mvid) WHERE m.mvname=(?)", key)
	CheckErr(err)
	for sqlquery.Next() {
		err = sqlquery.Scan(&getmov.Name, &getmov.AVG, &getmov.Count, &getmov.Comments)
	}

	out, err := json.Marshal(getmov)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, string(out))

}
