package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
)

var db *sql.DB

func uniqueID() string {
	guid := xid.New()
	return guid.String()
}

func createProfile(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(uniqueID()))
}

func createListing(res http.ResponseWriter, req *http.Request) {
	firstN := req.FormValue("firstName")
	lastN := req.FormValue("lastName")
	age := req.FormValue("age")
	gender := req.FormValue("gender")
	city := req.FormValue("city")
	state := req.FormValue("state")
	startDate := req.FormValue("startDate")
	endDate := req.FormValue("endDate")
	interests := req.FormValue("interests")
	uuid := req.FormValue("uuid")

	_, err := db.Exec("INSERT INTO travelListings(firstN, lastN, age, gender, city, state, startDate, endDate, interests, uuid) VALUES(?,?,?,?,?,?,?,?,?,?)", firstN, lastN, age, gender, city, state, startDate, endDate, interests, uuid)

	if err != nil {
		fmt.Println(err)
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	}
	res.Write([]byte("Listing Created!"))
}

func readMyListings(res http.ResponseWriter, req *http.Request) {

	uuid := req.FormValue("uuid")
	rows, err := db.Query("SELECT * FROM travelListings WHERE uuid=?", uuid)
	if err != nil {
		res.Write([]byte("Error"))
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		res.Write([]byte("Error"))
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		res.Write([]byte("Error"))
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func deleteListing(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	uuid := req.URL.Query().Get("UUID")
	query, err := db.Exec("DELETE FROM travelListings WHERE id=? AND uuid=?", id, uuid)
	if err != nil {
		http.Error(res, "Unable to delete book", 500)
		return
	}
	affected, err := query.RowsAffected()
	if err != nil {
		http.Error(res, "Server error", 500)
		return
	}
	if affected == 0 {
		res.Write([]byte("Failed to delete listing"))
		return
	}
	res.Write([]byte("Listing Successfully Deleted"))

}

func readMatchedListings(res http.ResponseWriter, req *http.Request) {
	sDate := req.FormValue("startDate")
	eDate := req.FormValue("endDate")
	city := req.FormValue("city")
	state := req.FormValue("state")
	rows, err := db.Query("SELECT * FROM travelListings WHERE startDate >= ? AND endDate <= ? AND city=? AND state=?", sDate, eDate, city, state)
	if err != nil {
		res.Write([]byte("Error"))
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		res.Write([]byte("Error"))
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		res.Write([]byte("Error"))
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func main() {
	db, err := sql.Open("mysql", "root:mhacks@/main")
	if err != nil {

		panic(err.Error())

	}
	fmt.Println("BLAH")

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/createProfile", createProfile)
	http.HandleFunc("/createListing", createListing)
	http.HandleFunc("/deleteListing", deleteListing)
	http.HandleFunc("/readMyListings", readMyListings)
	http.HandleFunc("/readMatchedListings", readMatchedListings)
	http.ListenAndServe(":8080", nil)
}
