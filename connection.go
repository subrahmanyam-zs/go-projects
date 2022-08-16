package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlConfig struct {
	User     string
	Host     string
	Password string
	Port     string
	DbName   string
}

func Connection(mysql MySqlConfig) (*sql.DB, error) {
	constr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.DbName)

	db, err := sql.Open("mysql", constr)
	if err != nil {
		fmt.Println(err)
	}
	db.Exec(`create table if not exists employee(id uuid primary key,name text,dob text,city text,majors text,dId int)`)
	return db, err
}
