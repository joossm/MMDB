package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

const post, get = "POST", "GET"

func InitDatabase(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case get:
		// create mysql database if not exists
		_ = runQuery("CREATE DATABASE IF NOT EXISTS mmdb", "createTable")
		// create shema if not exists
		//_ = runQuery("CREATE SCHEMA IF NOT EXISTS mmdb", "createTable")
		// create table if not exists for images with id, name, image
		_ = runQuery("CREATE TABLE IF NOT EXISTS images (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), image MEDIUMBLOB)", "createTable")
		// create table if not exists for user with id, username, password
		_ = runQuery("CREATE TABLE `users` ( `idusers` INT NOT NULL, `username` VARCHAR(45) NULL, `password` VARCHAR(45) NULL, PRIMARY KEY (`idusers`));", "createTable")
		// create table if not exists for generes with id, name
		_ = runQuery("CREATE TABLE IF NOT EXISTS generes (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255))", "createTable")
		// create table if not exists for combine of images and generes with id, image_id, genre_id
		_ = runQuery("CREATE TABLE IF NOT EXISTS images_generes (id INT AUTO_INCREMENT PRIMARY KEY, image_id INT, genre_id INT)", "createTable")
		// create table if not exists for combine of images and users with id, image_id, user_id
		_ = runQuery("CREATE TABLE IF NOT EXISTS images_users (id INT AUTO_INCREMENT PRIMARY KEY, image_id INT, user_id INT)", "createTable")

		_, responseErr := responseWriter.Write([]byte("Database initialized"))
		errorHandler(responseErr)
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
	switch function {
	case "createTable": // CREATE TABLE
		_, err := db.Exec(query)
		errorHandler(err)
		defer closeDB(db)
		return nil
	case "insert": // INSERT

	case "select": // SELECT
		result, err := db.Query(query)
		errorHandler(err)
		defer closeDB(db)
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
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	// return hello world
	js, err := json.Marshal("Hello World")
	errorHandler(err)
	_, responseErr := w.Write(js)
	errorHandler(responseErr)
}
func openDB() *sql.DB {
	fmt.Println("Opening DB")
	db, err := sql.Open("mysql", "admin:password@tcp(mysql:3306)/mmdb")
	fmt.Println(db.Ping())
	fmt.Println(db.Stats())
	db.SetMaxIdleConns(0)
	errorHandler(err)
	//defer closeDB(db)
	return db
}

func errorHandler(err error) {
	if err != nil {
		//panic(err) is not required, but it is a good idea to panic on error
		fmt.Println(err)
	}
}
