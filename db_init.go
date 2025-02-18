package main

import "database/sql"

var db *sql.DB

func initdatabase() error {
	dns := "root:fzfz1314@tcp(127.0.0.1:3306)/e_commerce?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = sql.Open("mysql", dns)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}
