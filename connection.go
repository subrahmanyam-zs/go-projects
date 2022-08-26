package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	User     string
	Host     string
	Password string
	Port     string
	DbName   string
}

func Connection(mysql *MySQLConfig) (*sql.DB, error) {
	constr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.DbName)

	db, err := sql.Open("mysql", constr)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`create table if not exists employee(id varchar(50) primary key,name varchar(20),
    dob varchar(10),city varchar(20),majors varchar(20),dId int)`)

	if err != nil {
		return nil, err
	}

	return db, err
}
