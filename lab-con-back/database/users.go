package database

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	uid      int
	username string
	phone    int
	email    string
	level    int
	schoolId int
	created  int64
	password string
}

func NewUser(username, password, email string, schoolId, level, phone int) User {
	u := User{
		username: username,
		level:    level,
		email:    email,
		phone:    phone,
		schoolId: schoolId,
		password: password,
		created:  time.Now().Unix(),
	}
	return u
}

func (d DB) InitialUser() {
	var db = d.GetHandler()
	var err error
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS user(
		uid INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(128) NOT NULL ,
		level INTEGER DEFAULT 0,
		phone INTEGER NULL ,
		email VARCHAR(128) NULL,
		schoolId INTEGER NOT NULL ,
		password VARCHAR(128) NOT NULL ,
        created DATE NULL);`)
	if err != nil {
		panic(err)
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}
}
func (d DB) InsertUser(u User) (int64, error) {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if db.QueryRow(`select * from user where username=? or schoolId=?`, u.username, u.schoolId).Scan() != nil {
		return 0, errors.New("user Already Exist")
	}
	var err error
	exec, err := db.Exec(`INSERT INTO user(username,level,password,schoolId,created,email,phone) VALUES (?,?,?,?,?,?,?)`,
		u.username, u.level, u.password, u.schoolId, u.created, u.email, u.phone)
	if err != nil {
		panic(err)
	}
	id, err := exec.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id, nil
}
func (d DB) checkPassword(schoolId int, password string) bool {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	rows, _ := db.Query("select * from user where schoolId=? and password=?", schoolId, password)
	return rows.Next()
}
func (d DB) Login(u User, password string) (bool, error) {
	var db = d.GetHandler()
	if u.email != "" {
		var schoolId int
		row := db.QueryRow("SELECT schoolId from user where email=?", u.email)
		if row.Scan(&schoolId) != nil {
			return false, errors.New("no such user with that email")
		} else {
			return d.checkPassword(schoolId, password), nil
		}
	}
	if u.phone != 0 {
		var schoolId int
		row := db.QueryRow("SELECT schoolId from user where phone=?", u.phone)
		if row.Scan(&schoolId) != nil {
			return false, errors.New("no such user with that phone")
		} else {
			return d.checkPassword(schoolId, password), nil
		}
	}
	if u.schoolId != 0 {
		var schoolId int
		row := db.QueryRow("SELECT schoolId from user where schoolId=?", u.schoolId)
		if row.Scan(&schoolId) != nil {
			return false, errors.New("no such user with that schoolId")
		} else {
			return d.checkPassword(schoolId, password), nil
		}
	}
	return false, errors.New("user info is empty")
}
func (d DB) UpdateEmail(schoolId int, email string) error {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	var uid int
	row := db.QueryRow("select uid from user where schoolId=?", schoolId)
	if row.Scan(&uid) != nil {
		return errors.New("no such user with that schoolId")
	}
	_, err := db.Exec("update user set email=? where uid=?", email, uid)
	return err
}
func (d DB) UpdatePhone(schoolId, phone int) error {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	var uid int
	row := db.QueryRow("select uid from user where schoolId=?", schoolId)
	if row.Scan(&uid) != nil {
		return errors.New("no such user with that schoolId")
	}
	_, err := db.Exec("update user set phone=? where uid=?", phone, uid)
	return err
}
func (d DB) UpdatePassword(schoolId int, password string) error {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	var uid int
	row := db.QueryRow("select uid from user where schoolId=?", schoolId)
	if row.Scan(&uid) != nil {
		return errors.New("no such user with that schoolId")
	}
	_, err := db.Exec("update user set password=? where uid=?", password, uid)
	return err
}
func (d DB) CheckPermission(schoolId, required int) (bool, error) {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	var level int
	row := db.QueryRow("select level from user where schoolId=?", schoolId)
	if row.Scan(&level) != nil {
		return false, errors.New("no such user with that schoolId")
	} else {
		return level > required, nil
	}
}
