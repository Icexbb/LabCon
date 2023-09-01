package database

import (
	"database/sql"
	"errors"
	"time"
)

type Lab struct {
	id          int
	name        string
	description string
	created     int64
}

func NewLab(name, des string) Lab {
	var lab = Lab{
		name:        name,
		description: des,
		created:     time.Now().Unix(),
	}
	return lab
}

func (d DB) InitialLab() {
	var db = d.GetHandler()
	var err error
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS lab(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(128) NOT NULL ,
		description INTEGER DEFAULT 0,
        created DATE NULL);`)
	if err != nil {
		panic(err)
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}
}
func (d DB) InsertLab(l Lab) (int64, error) {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if db.QueryRow("select * from lab where name=?", l.name).Scan() != nil {
		return 0, errors.New("lab name already exist")
	}
	var err error
	exec, err := db.Exec(`INSERT INTO lab(name,description,created) VALUES (?,?,?)`,
		l.name, l.description, l.created)
	if err != nil {
		panic(err)
	}
	id, err := exec.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id, nil
}

type LabAppRecord struct {
	LabId      int   `json:"lab_id"`
	UserId     int   `json:"user_id"`
	Semester   int   `json:"semester"`
	WeekFrom   int   `json:"week_from"`
	WeekTO     int   `json:"week_to"`
	Weekday    int   `json:"weekday"`
	PeriodFrom int   `json:"period_from"`
	PeriodTo   int   `json:"period_to"`
	Created    int64 `json:"created"`
}

func newRecord(userId, labId, semester, weekFrom, weekTO, weekday, periodFrom, periodTo int) LabAppRecord {
	return LabAppRecord{
		UserId:     userId,
		LabId:      labId,
		Semester:   semester,
		WeekFrom:   weekFrom,
		WeekTO:     weekTO,
		Weekday:    weekday,
		PeriodFrom: periodFrom,
		PeriodTo:   periodTo,
		Created:    time.Now().Unix(),
	}
}

func (d DB) InitialRecord() {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	_, _ = db.Exec(
		`CREATE TABLE IF NOT EXISTS record(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		userId     INTEGER not null ,
		labId      INTEGER not null ,
		semester   INTEGER not null ,
		weekFrom   INTEGER not null ,
		weekTO     INTEGER not null ,
		weekday    INTEGER not null ,
		periodFrom INTEGER not null ,
		periodTo   INTEGER not null ,
        created DATE NULL);`)
}
func (d DB) GetLabRecord(labId int) []LabAppRecord {
	var result []LabAppRecord
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	rows, err := db.Query(`select * from record where labId=?`, labId)
	if err != nil {
		return result
	}
	for rows.Next() {
		var lid int
		var userId int
		var semester int
		var weekFrom int
		var weekTO int
		var weekday int
		var periodFrom int
		var periodTo int
		var created time.Time
		err := rows.Scan(&lid, &userId, &semester, &weekFrom, &weekTO, &weekday, &periodFrom, &periodTo, &created)
		if err != nil {
			return nil
		}
		var r = LabAppRecord{
			LabId:      lid,
			UserId:     userId,
			Semester:   semester,
			WeekFrom:   weekFrom,
			WeekTO:     weekTO,
			Weekday:    weekday,
			PeriodFrom: periodFrom,
			PeriodTo:   periodTo,
			Created:    created.Unix(),
		}
		result = append(result, r)
	}
	return result
}
