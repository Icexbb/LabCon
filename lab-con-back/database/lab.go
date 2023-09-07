package database

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Lab struct {
	Id                         int
	Name, Description, Created string
}

func NewLab(name, des string) Lab {
	var lab = Lab{
		Name:        name,
		Description: des,
		Created:     time.Now().Format(time.RFC3339),
	}
	return lab
}

func (d DB) InitialLab() {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS lab(
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(128) NOT NULL ,
		description VARCHAR(1024) DEFAULT 0,
        created TIMESTAMP NULL);`)
	if err != nil {
		panic(err)
	}

}
func (d DB) InsertLab(l Lab) (int64, error) {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if db.QueryRow("select * from lab where name=?", l.Name).Scan() == nil {
		return -1, errors.New("lab name already exist")
	}
	exec, err := db.Exec(
		`INSERT INTO lab(name,description,created) VALUES (?,?,?)`, l.Name, l.Description, l.Created)
	if err != nil {
		return -1, err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		log.Fatalf("InsertLab: %v", err)
		return -1, err
	}
	return id, nil
}

func (d DB) GetAllLabs() ([]Lab, error) {
	db := d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	result := []Lab{}
	rows, err := db.Query(`SELECT id,name,description,created FROM lab;`)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var id int
		var name, description string
		var created time.Time
		err = rows.Scan(&id, &name, &description, &created)
		if err != nil {
			return result, err
		}
		r := Lab{
			Id:          id,
			Name:        name,
			Description: description,
			Created:     created.Format(time.RFC3339),
		}
		result = append(result, r)
	}
	return result, nil
}
func (d DB) GetLabById(labId int) (Lab, error) {
	db := d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	var name, description string
	var created time.Time
	err := db.QueryRow(`SELECT name,description,created FROM lab WHERE id=?;`, labId).Scan(&name, &description, &created)
	if err != nil {
		return Lab{}, err
	} else {
		return Lab{labId, name, description, created.Format(time.RFC3339)}, nil
	}
}
func (d DB) GetLabByName(name string) (Lab, error) {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	var labId int
	var description string
	var created time.Time
	err := db.QueryRow(`SELECT id,description,created FROM lab WHERE name=?;`, name).Scan(&labId, &description, &created)
	if err != nil {
		return Lab{}, err
	} else {
		return Lab{labId, name, description, created.Format(time.RFC3339)}, err
	}
}

func (d DB) DeleteLabByName(name string) (int64, error) {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if db.QueryRow("select * from lab where name=?", name).Scan() == nil {
		return 0, errors.New("lab name not exist")
	}
	exec, err := db.Exec("DELETE from lab where name=?", name)
	affected, err := exec.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, err
}
func (d DB) DeleteLabById(labId int) (int64, error) {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if db.QueryRow("select * from lab where id=?", labId).Scan() == nil {
		return 0, errors.New("lab ID not exist")
	}
	exec, err := db.Exec("DELETE from lab where id=?", labId)
	affected, err := exec.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, err
}

func (d DB) UpdateLabName(labId int, name string) error {
	_, err := d.GetLabById(labId)
	if err == nil {
		db := d.GetHandler()
		defer func(db *sql.DB) {
			_ = db.Close()
		}(db)
		_, err = db.Exec(`UPDATE lab set name=? where id=?`, name, labId)
	}
	return err
}
func (d DB) UpdateLabDescription(labId int, des string) error {
	_, err := d.GetLabById(labId)
	if err == nil {
		db := d.GetHandler()
		defer func(db *sql.DB) {
			_ = db.Close()
		}(db)
		_, err = db.Exec(`UPDATE lab set description=? where id=?`, des, labId)
	}
	return err
}

type LabAppRecord struct {
	LabId       int
	UserId      int
	Semester    int
	WeekFrom    int
	WeekTo      int
	Weekday     int
	PeriodFrom  int
	PeriodTo    int
	Description string
	Created     string
}

func NewRecord(userId, labId, semester, weekFrom, weekTO, weekday, periodFrom, periodTo int, Description string) LabAppRecord {
	return LabAppRecord{
		UserId:      userId,
		LabId:       labId,
		Semester:    semester,
		WeekFrom:    weekFrom,
		WeekTo:      weekTO,
		Weekday:     weekday,
		PeriodFrom:  periodFrom,
		PeriodTo:    periodTo,
		Description: Description,
		Created:     time.Now().Format(time.RFC3339),
	}
}

func (d DB) InitialRecord() {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	_, _ = db.Exec(
		`CREATE TABLE IF NOT EXISTS record(
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		userId INTEGER not null ,
		labId INTEGER not null ,
		Semester INTEGER not null ,
		weekFrom INTEGER not null ,
		weekTo INTEGER not null ,
		weekday INTEGER not null ,
		periodFrom INTEGER not null ,
		periodTo INTEGER not null ,
		description VARCHAR(128) default null,
        created TIMESTAMP NULL);`)
}
func (d DB) InsertLabRecord(record LabAppRecord) (int64, error) {
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	exec, err := db.Exec(
		`INSERT INTO record(
                        userId,labId,Semester,weekFrom,weekTo,weekday,periodFrom,periodTo,description,created
                        ) VALUES (?,?,?,?,?,?,?,?,?,?)`,
		record.UserId, record.LabId, record.Semester, record.WeekFrom, record.WeekTo, record.Weekday,
		record.PeriodFrom, record.PeriodTo, record.Description, record.Created)
	if err != nil {
		return -1, err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
func (d DB) GetLabRecord(labId int) []LabAppRecord {
	var result []LabAppRecord
	var db = d.GetHandler()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	rows, err := db.Query(
		`SELECT userId,Semester,weekFrom,weekTo,weekday,periodFrom,periodTo,description,created FROM lab.record WHERE labId=?;`, labId)
	if err != nil {
		return result
	}

	for rows.Next() {
		var userId, semester, weekFrom, weekTo, weekday, periodFrom, periodTo int
		var description string
		var created time.Time
		err = rows.Scan(&userId, &semester, &weekFrom, &weekTo, &weekday, &periodFrom, &periodTo, &description, &created)
		if err != nil {
			return nil
		}
		var r = LabAppRecord{
			LabId: labId, UserId: userId, Semester: semester,
			WeekFrom: weekFrom, WeekTo: weekTo, Weekday: weekday,
			PeriodFrom: periodFrom, PeriodTo: periodTo,
			Description: description,
			Created:     created.Format(time.RFC3339),
		}

		result = append(result, r)
	}
	return result
}
func (d DB) CheckFreePeriod(labId, semesterId, weekNum int) [7][]int {
	s := d.GetSemesterInfo(semesterId)
	var result [7][]int

	{
		var dayCount []int
		for i := 1; i <= s.periodCount; i++ {
			dayCount = append(dayCount, i)
		}
		for j := 0; j < 7; j++ {
			result[j] = dayCount
		}
	}

	records := d.GetLabRecord(labId)

	for _, record := range records {
		if record.WeekFrom <= weekNum && weekNum <= record.WeekTo {
			dayCount := result[record.Weekday-1]
			var dayResult []int
			for _, i := range dayCount {
				if record.PeriodFrom > i || i > record.PeriodTo {
					dayResult = append(dayResult, i)
				}
			}
			result[record.Weekday-1] = dayResult
		}
	}
	return result
}
