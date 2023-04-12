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
		_ = runQuery("CREATE TABLE IF NOT EXISTS image (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), image MEDIUMBLOB)", "createTable")
		// create table if not exists for user with id, username, password
		_ = runQuery("CREATE TABLE `user` ( `idusers` INT NOT NULL, `username` VARCHAR(45) NULL, `password` VARCHAR(45) NULL, PRIMARY KEY (`idusers`));", "createTable")
		// create table if not exists for genres with id, name
		_ = runQuery("CREATE TABLE IF NOT EXISTS genre (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255))", "createTable")
		// create table if not exists for combine of images and genres with id, image_id, genre_id
		_ = runQuery("CREATE TABLE IF NOT EXISTS image_genre (id INT AUTO_INCREMENT PRIMARY KEY, image_id INT, genre_id INT)", "createTable")
		// create table if not exists for combine of images and users with id, image_id, user_id
		_ = runQuery("CREATE TABLE IF NOT EXISTS image_user (id INT AUTO_INCREMENT PRIMARY KEY, image_id INT, user_id INT)", "createTable")

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
		_, err := db.Exec(query, args...)
		errorHandler(err)
		defer closeDB(db)
		return nil
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

type Image struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image []byte `json:"image"`
}

func Index(responseWriter http.ResponseWriter, request *http.Request) {
	// Load images from database and display them on the index page
	switch request.Method {
	case get:
		// get all images from database
		rows := runQuery("SELECT * FROM images", "select")
		// create array of images
		var images []Image
		// loop through all images
		for rows.Next() {
			// create image
			var image Image
			// scan image
			err := rows.Scan(&image.ID, &image.Name, &image.Image)
			errorHandler(err)
			// append image to array
			images = append(images, image)
		}
		// return images
		js, err := json.Marshal(images)
		errorHandler(err)
		_, responseErr := responseWriter.Write(js)
		errorHandler(responseErr)
		return
	}

}

func UploadImage(responseWriter http.ResponseWriter, request *http.Request) {

}
func DownloadImage(responseWriter http.ResponseWriter, request *http.Request) {

}
func DeleteImage(responseWriter http.ResponseWriter, request *http.Request) {

}
func RegisterUser(responseWriter http.ResponseWriter, request *http.Request) {
	// check request method
	switch request.Method {
	case get:
		// return register.html
		http.ServeFile(responseWriter, request, "./view/register.html")
	case post:
		// get username and password from request
		username := request.FormValue("username")
		password := request.FormValue("password")
		// check if username and password are not empty
		if username != "" && password != "" {
			// Check if username already exists
			rows := runQuery("SELECT * FROM users WHERE username = ?", "select", username)
			// if username does not exist
			if !rows.Next() {
				// register user
				_ = runQuery("INSERT INTO users (username, password) VALUES (?, ?)", "insert", username, password)
				// return success
				js, err := json.Marshal("User registered")
				errorHandler(err)
				_, responseErr := responseWriter.Write(js)
				errorHandler(responseErr)
			}
		}
	}
}
func LoginUser(responseWriter http.ResponseWriter, request *http.Request) {
	// check request method
	switch request.Method {
	case get:
		// return login.html
		http.ServeFile(responseWriter, request, "./view/login.html")
	case post:
		// get username and password from request
		username := request.FormValue("username")
		password := request.FormValue("password")
		// check if username and password are not empty
		if username != "" && password != "" {
			// Check if username and password are correct
			rows := runQuery("SELECT * FROM users WHERE username = ? AND password = ?", "select", username, password)
			// if correct
			if rows.Next() {
				// login user
				// return success
				js, err := json.Marshal("User logged in")
				errorHandler(err)
				_, responseErr := responseWriter.Write(js)
				errorHandler(responseErr)
			} else {
				// return error
				js, err := json.Marshal("Username or password is incorrect")
				errorHandler(err)
				_, responseErr := responseWriter.Write(js)
				errorHandler(responseErr)
			}
		}
	}

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
