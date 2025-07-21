package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	var (
		err error
		dsn = "root:root@tcp(127.0.0.1:3306)/test01?charset=utf8mb4&parseTime=true&loc=Local"
	)

	// 1. 使用 sqlx.Open 连接数据库
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 使用 sqlx.Open 变体方法 sqlx.MustOpen 连接数据库，如果出现错误直接 panic
	db = sqlx.MustOpen("mysql", dsn)

	// 3. 如果已经有了 *sql.DB 对象，则可以使用 sqlx.NewDb 连接数据库，得到 *sqlx.DB 对象
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db = sqlx.NewDb(sqlDB, "mysql")

	// 使用前 3 种方式连接数据库并不会立即与数据库建立连接，连接将会在合适的时候延迟建立。
	// 为了确保能够正常连接数据库，往往需要调用 db.Ping() 方法进行验证
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// 4. 使用 sqlx.Connect 连接数据库，等价于 sqlx.Open + db.Ping
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 5. 使用 sqlx.Connect 变体方法 sqlx.MustConnect 连接数据库，如果出现错误直接 panic
	db = sqlx.MustConnect("mysql", dsn)

	///////  MustXxx 类似方法名时就应该想到，其功能往往等价于 Xxx 方法，
	///////  不过在其内部实现中，遇到 error 不再返回，而是直接进行 panic，这也是 Go 语言很多库中的惯用方法

}
