package main

import (
	"log"

	"lab-con-front/database"
)

func main() {
	var dbPath = "./db.sqlite"

	db := database.NewDB("sqlite3", dbPath)
	db.InitialLab()
	db.InitialUser()
	u := database.NewUser("admin", "password", "", 10000000000000, 10, 0)
	userId, err := db.InsertUser(u)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(userId)
	}
	logResult, err := db.Login(u, "password1")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(logResult)
	}
}
