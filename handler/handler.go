package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

const post, get = "POST", "GET"

func InitDatabase(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case get:
		// create table if not exists for images with id, name, image
		_ = runQuery("CREATE TABLE IF NOT EXISTS images (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), image MEDIUMBLOB)", "createTable")
		// create table if not exists for user with id, username, password
		_ = runQuery("CREATE TABLE IF NOT EXISTS users (id INT AUTO_INCREMENT PRIMARY KEY, username VARCHAR(255), password VARCHAR(255))", "createTable")

		return
	default:
		js, err := json.Marshal("THIS IS A GET REQUEST")
		errorHandler(err)
		_, responseErr := responseWriter.Write(js)
		errorHandler(responseErr)
		return
	}
}
func runQuery(query string, function string, args ...interface{}) *sql.Rows {
	db := openDB()
	defer closeDB(db)
	switch function {
	case "createTable": // CREATE TABLE
		_, err := db.Exec(query)
		errorHandler(err)
		return nil
	case "insert": // INSERT

	case "select": // SELECT
		result, err := db.Query(query)
		errorHandler(err)
		return result
	case "update": // UPDATE

	case "delete": // DELETE

	}
	return nil
}

func InsertMusik(w http.ResponseWriter, r *http.Request) {

}
func closeDB(db *sql.DB) {
	err := db.Close()
	errorHandler(err)
}

func openDB() *sql.DB {
	fmt.Println("Opening DB")
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/books")
	fmt.Println(db.Ping())
	fmt.Println(db.Stats())
	db.SetMaxIdleConns(0)
	errorHandler(err)
	defer closeDB(db)
	return db
}

func errorHandler(err error) {
	if err != nil {
		//panic(err) is not required, but it is a good idea to panic on error
		fmt.Println(err)
	}
}
