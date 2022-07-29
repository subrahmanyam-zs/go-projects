package driver

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func ConnectToSQL() (*sql.DB, error) {
	configure := mysql.Config{
		User:                 "root",
		Passwd:               "Bansal@8023",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "test",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", configure.FormatDSN())
	if err != nil {
		log.Println(err)

		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println(err)

		return nil, err
	}

	log.Println("Connected")

	return db, nil
}
