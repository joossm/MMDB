package handler

import (
	"MMDB/model"
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
		// create table if not exists for images with id, name, image
		_ = runQuery("CREATE TABLE IF NOT EXISTS image (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), image MEDIUMBLOB)", "createTable")
		// create table if not exists for user with id, username, password
		_ = runQuery("CREATE TABLE IF NOT EXISTS user (id INT AUTO_INCREMENT PRIMARY KEY, username VARCHAR(255), password VARCHAR(255))", "createTable")
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
		rows := runQuery("SELECT * FROM image", "select")
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
	// request type switch
	switch request.Method {
	case get:
		// return upload.html site
		http.ServeFile(responseWriter, request, "./view/uploadImage.html")
		return
	case post:
		// get image from request
		image, _, err := request.FormFile("image")
		errorHandler(err)
		// get name from request
		name := request.FormValue("name")
		// save image to database
		_ = runQuery("INSERT INTO image (name, image) VALUES (?, ?);", "insert", name, image)
		// return index.html site
		http.ServeFile(responseWriter, request, "./view/index.html")
		return
	}
}
func DownloadImage(responseWriter http.ResponseWriter, request *http.Request) {

}
func DeleteImage(responseWriter http.ResponseWriter, request *http.Request) {

}
func Register(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case get: // GET
		// return register.html site
		http.ServeFile(responseWriter, request, "./view/register.html")

	case post:
		fmt.Println("Register was executed")

		user := model.User{}
		user.Username = request.FormValue("name")
		user.Password = request.FormValue("password")

		db := openDB()
		defer closeDB(db)
		result, err := db.Query("SELECT username FROM user WHERE username = ?", user.Username)
		fmt.Println("result: ", result)
		errorHandler(err)
		fmt.Println("Query executed")
		var users []model.User
		if result.Next() == true {
			for result.Next() {
				var user model.User
				err = result.Scan(&user.Id, &user.Username, &user.Password)
				fmt.Println("user: ", user.Username, user.Password)
				users = append(users, user)
			}
			if users != nil {
				js, err := json.Marshal("already exists")
				errorHandler(err)
				_, responseErr := responseWriter.Write(js)
				errorHandler(responseErr)
				return
			}
		} else {
			// GET MAX ID
			result, err := db.Query("SELECT MAX(id) FROM user")
			errorHandler(err)
			var maxId int
			if result != nil {
				for result.Next() {
					err = result.Scan(&maxId)
					errorHandler(err)
				}
			}
			maxId++
			fmt.Println("result is nil | execute insert")
			res, err := db.Query("INSERT INTO user (id, username, password) VALUES (?,?,?)",
				maxId, user.Username, user.Password)
			fmt.Println(res)
			errorHandler(err)
			js, err := json.Marshal("true")
			_, responseErr := responseWriter.Write(js)
			errorHandler(responseErr)
			return
		}

	default:
		js, err := json.Marshal("THIS IS A POST REQUEST")
		errorHandler(err)
		_, responseErr := responseWriter.Write(js)
		errorHandler(responseErr)
		return
	}
}

func Login(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case get:
		// display login.html
		http.ServeFile(responseWriter, request, "./view/login.html")
	case post:
		user := model.User{}
		user.Username = request.FormValue("name")
		user.Password = request.FormValue("password")

		db := openDB()
		defer closeDB(db)
		result, err := db.Query("SELECT * FROM user WHERE username = ? AND password = ?", user.Username, user.Password)
		errorHandler(err)
		var users []model.User
		if result != nil {
			for result.Next() {
				var user model.User
				err = result.Scan(&user.Id, &user.Username, &user.Password)
				errorHandler(err)
				users = append(users, user)
			}
			// redirect to index page
			http.ServeFile(responseWriter, request, "./view/index.html")
			return
		} else {
			// return error
			js, err := json.Marshal("false")
			errorHandler(err)
			_, responseErr := responseWriter.Write(js)
			errorHandler(responseErr)
			// redirect to login page
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
			return
		}
	default:
		js, err := json.Marshal("THIS IS A POST REQUEST")
		errorHandler(err)
		_, responseErr := responseWriter.Write(js)
		errorHandler(responseErr)
		return

	}
}
func closeDB(db *sql.DB) {
	err := db.Close()
	errorHandler(err)
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
