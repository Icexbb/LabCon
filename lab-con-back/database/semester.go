package database

import (
	"database/sql"
	"errors"
	"time"
)

type Semester struct {
	id          int
	name        string
	description string
	weekCount   int
	periodCount int
	start       string
	end         string
	created     string
}

func NewSemester(name, description string, weekCount, periodCount int, start, end time.Time) Semester {
	s := Semester{
		name:        name,
		description: description,
		weekCount:   weekCount,
		periodCount: periodCount,
		start:       start.Format(time.RFC3339),
		end:         end.Format(time.RFC3339),
		created:     time.Now().Format(time.RFC3339),
	}
	return s
}
func (d DB) InitialSemester() {
	db := d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS semester(
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(128) NOT NULL,
		description VARCHAR(128) DEFAULT null,
		weekCount integer not null,
		periodCount integer not null,
		start TIMESTAMP not null,
		end TIMESTAMP not null,
		created TIMESTAMP NULL);`)
	if err != nil {
		panic(err)
	}
}
func (d DB) InsertSemester(s Semester) (int64, error) {
	db := d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if db.QueryRow(`select * from semester where name=?`, s.name).Scan() == nil {
		return -1, errors.New("semester name already exist")
	}
	exec, err := db.Exec(
		`INSERT INTO semester(name,description,weekCount,periodCount,start,end,created) VALUES (?,?,?,?,?,?,?)`,
		s.name, s.description, s.weekCount, s.periodCount, s.start, s.end, s.created)
	if err != nil {
		return -1, err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
func (d DB) GetSemesterInfo(id int) Semester {
	db := d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	rows := db.QueryRow(`select * from semester where id=?`, id)
	var sid int
	var name, description string
	var weekCount, periodCount int
	var start, end, created string
	err := rows.Scan(&sid, &name, &description, &weekCount, &periodCount, &start, &end, &created)
	if err != nil {
		return Semester{}
	} else {
		return Semester{
			sid, name, description,
			weekCount, periodCount, start, end, created,
		}
	}
}
