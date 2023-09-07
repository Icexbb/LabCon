package database

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	Uid, Phone, Level, SchoolId        int
	Username, Email, Created, Password string
}

func NewUser(username, password, email string, schoolId, level, phone int) User {
	u := User{
		Username: username,
		Level:    level,
		Email:    email,
		Phone:    phone,
		SchoolId: schoolId,
		Password: password,
		Created:  time.Now().Format(time.RFC3339),
	}
	return u
}

func (d DB) InitialUser() {
	var db = d.GetHandler()
	var err error
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS user(
		uid INTEGER PRIMARY KEY AUTO_INCREMENT,
		username VARCHAR(128) NOT NULL ,
		level INTEGER DEFAULT 0,
		phone INTEGER NULL ,
		email VARCHAR(128) NULL,
		schoolId BIGINT NOT NULL ,
		password VARCHAR(128) NOT NULL ,
        created TIMESTAMP NULL);`)
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
	if !errors.Is(db.QueryRow(`SELECT * FROM user WHERE username=? or schoolId=?`, u.Username, u.SchoolId).Scan(), sql.ErrNoRows) {
		return 0, errors.New("user Already Exist")
	}
	var err error
	exec, err := db.Exec(`INSERT INTO user(username,level,password,schoolId,created,email,phone) VALUES (?,?,?,?,?,?,?)`,
		u.Username, u.Level, u.Password, u.SchoolId, u.Created, u.Email, u.Phone)
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
	rows, _ := db.Query("select * from user where schoolId=? and password=?", schoolId, password)
	result := rows.Next()
	_ = db.Close()
	return result
}
func (d DB) Login(u User, password string) (bool, error) {
	var db = d.GetHandler()

	var schoolId int
	if u.SchoolId != 0 {
		row := db.QueryRow("SELECT schoolId from user where schoolId=?", u.SchoolId)
		if row.Scan(&schoolId) != nil {
			return false, errors.New("no such user with that schoolId")
		}
	} else if u.Email != "" {
		row := db.QueryRow("SELECT schoolId from user where email=?", u.Email)
		if row.Scan(&schoolId) != nil {
			return false, errors.New("no such user with that email")
		}
	} else if u.Phone != 0 {
		row := db.QueryRow("SELECT schoolId from user where phone=?", u.Phone)
		if row.Scan(&schoolId) != nil {
			return false, errors.New("no such user with that phone")
		}
	}
	_ = db.Close()
	if schoolId != 0 {
		return d.checkPassword(schoolId, password), nil
	} else {
		return false, errors.New("user info is empty")
	}
}
func (d DB) FindUserById(id int) User {
	result := User{}
	return result

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
func (d DB) UpdatePhone(schoolId int, phone int) error {
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
