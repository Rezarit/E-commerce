package dao

import "database/sql"

var DB *sql.DB

func Initdatabase() error {
	dns := "root:fzfz1314@tcp(127.0.0.1:3306)/e_commerce?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = sql.Open("mysql", dns)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}
	return nil
}
