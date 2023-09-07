package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lab-con-front/database"
)

func serve() {
	server := gin.Default()
	err := server.Run()
	if err != nil {
		return
	}
}
func InitialDatabase(dbType, dbPath string) database.DB {
	db := database.NewDB(dbType, dbPath)
	db.InitialLab()
	db.InitialUser()
	db.InitialRecord()
	db.InitialSemester()
	return db
}
func TestDatabase(db database.DB) {
	//u := database.NewUser(
	//	"admin", "password", "", 10000000000000, 10, 0)

	//userId, err := db.InsertUser(u)
	//if err != nil {
	//	log.Println("InsertUser:", err)
	//} else {
	//	log.Println("InsertUser:", userId)
	//}
	//logResult, err := db.Login(u, "password")
	//if err != nil {
	//	log.Println("Login:", err)
	//} else {
	//	log.Println("Login:", logResult)
	//}

	//lab := database.NewLab("test", "test lab")
	//labId, err := db.InsertLab(lab)
	//if err != nil {
	//	log.Println("InsertLab:", err)
	//} else {
	//	log.Println("InsertLab:", labId)
	//}

	//sid, err := db.InsertSemester(database.NewSemester("test", "", 16, 12,
	//	time.Date(2023, 9, 11, 0, 0, 0, 0, time.Local),
	//	time.Date(2024, 1, 28, 0, 0, 0, 0, time.Local)))
	//if err != nil {
	//	log.Println("InsertSemester:", err)
	//} else {
	//	log.Println("InsertSemester:", sid)
	//}
	//userId := 1
	labId := 4
	sid := 1
	//_, _ = db.InsertLabRecord(database.NewRecord(int(userId), int(labId), int(sid), 1, 10, 1, 4, 6, ""))
	//_, _ = db.InsertLabRecord(database.NewRecord(int(userId), int(labId), int(sid), 1, 10, 3, 1, 4, ""))
	free := db.CheckFreePeriod(int(labId), int(sid), 1)
	fmt.Printf("%v", free)
	//affected, err := db.DeleteLabByName("test")
	//if err != nil {
	//	log.Println("DeleteLabByName:", err)
	//} else {
	//	log.Println("DeleteLabByName:", affected)
	//}
}

func main() {
	//var dbPath = "./db.sqlite"

	dbPath := "lab:root@tcp(127.0.0.1:3306)/lab?charset=utf8&parseTime=True"
	db := InitialDatabase("mysql", dbPath)
	TestDatabase(db)
}
