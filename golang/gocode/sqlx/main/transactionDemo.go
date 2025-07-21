package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// 事务
func Transaction(db *sqlx.DB, id int64, name string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("UPDATE user SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}

	// 获取执行后影响的行数
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("rowsAffected: %d\n", rowsAffected)

	// 一旦 Commit 成功，之前 defer 的 Rollback 仍然会执行，但它会因为事务已提交而失败并返回一个可被忽略的错误，并不影响程序的正常执行
	return tx.Commit()
}

// 预处理
func PreparexGetUser(db *sqlx.DB) (User, error) {
	stmt, err := db.Preparex(`SELECT * FROM user WHERE id = ?`)
	if err != nil {
		return User{}, err
	}

	var u User
	err = stmt.Get(&u, 1)
	return u, err
}
