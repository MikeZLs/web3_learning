package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// QueryxUsers 查询多条记录
func QueryxUsers(db *sqlx.DB) ([]User, error) {
	var us []User
	rows, err := db.Queryx("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		// sqlx 提供了便捷方法可以将查询结果直接扫描到结构体
		err = rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}
	return us, nil
}

// QueryRowxUser 查询单条记录
func QueryRowxUser(db *sqlx.DB, id int) (User, error) {
	var u User
	err := db.QueryRowx("SELECT * FROM user WHERE id = ?", id).StructScan(&u)
	return u, err
}

// *sqlx.DB.Get 方法包装了 *sqlx.DB.QueryRowx 方法，用以简化查询单条记录
func GetUser(db *sqlx.DB, id int) (User, error) {
	var u User
	// 查询记录扫描数据到 struct
	err := db.Get(&u, "SELECT * FROM user WHERE id = ?", id)
	return u, err
}

// *sqlx.DB.Select 方法包装了 *sqlx.DB.Queryx 方法，用以简化查询多条记录
func SelectUsers(db *sqlx.DB) ([]User, error) {
	var us []User
	// 查询记录扫描数据到 slice
	err := db.Select(&us, "SELECT * FROM user")
	return us, err
}

// sqlx.In 方法
func SqlxIn(db *sqlx.DB, ids []int64) ([]User, error) {
	query, args, err := sqlx.In("SELECT * FROM user WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	//query = db.Rebind(query)  // mysql 可不用
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var us []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age,
			&user.Birthday, &user.Salary, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		us = append(us, user)
	}
	return us, nil
}

func NamedExec(db *sqlx.DB) error {
	m := map[string]interface{}{
		"email": "jianghushinian007@outlook.com",
		"age":   18,
	}
	result, err := db.NamedExec(`UPDATE user SET age = :age WHERE email = :email`, m)
	if err != nil {
		return err
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func NamedQuery(db *sqlx.DB) ([]User, error) {
	u := User{
		Email: "jianghushinian007@outlook.com",
		Age:   18,
	}
	rows, err := db.NamedQuery("SELECT * FROM user WHERE email = :email OR age = :age", u)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// 将结果扫描到 map
func MapScan(db *sqlx.DB) ([]map[string]interface{}, error) {
	rows, err := db.Queryx("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []map[string]interface{}
	for rows.Next() {
		r := make(map[string]interface{})
		err := rows.MapScan(r)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, err
}

// 将结果扫描到 slice
func SliceScan(db *sqlx.DB) ([][]interface{}, error) {
	rows, err := db.Queryx("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res [][]interface{}
	for rows.Next() {
		// cols is an []interface{} of all the column results
		cols, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}
		res = append(res, cols)
	}
	return res, err
}

func main() {
	//users, err := QueryxUsers(db)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(users)

	users, err := SqlxIn(db, []int64{1, 2})
	if err != nil {
		panic(err)
	}
	fmt.Println(users)
}
