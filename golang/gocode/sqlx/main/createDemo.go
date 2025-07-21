package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

func MustCreateUser(db *sqlx.DB) (int64, error) {
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	user := User{
		Name:     sql.NullString{String: "jianghushinian", Valid: true},
		Email:    "jianghushinian007@outlook.com",
		Age:      10,
		Birthday: birthday,
		Salary: Salary{
			Month: 100000,
			Year:  10000000,
		},
	}

	res := db.MustExec(
		`INSERT INTO user(name, email, age, birthday, salary) VALUES(?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.Age, user.Birthday, user.Salary,
	)
	return res.LastInsertId()
}

//func main() {
//	//MustCreateUser(db)
//
//}
