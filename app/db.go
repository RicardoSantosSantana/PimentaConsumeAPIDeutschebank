package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func StringConnection() string {

	settings := Settings()

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.Database.User,
		settings.Database.Password,
		settings.Database.Host,
		settings.Database.Port,
		settings.Database.Name)

}

func openConnection() (*sql.DB, error) {

	strConn := StringConnection()

	db, err := sql.Open("mysql", strConn)
	db.SetConnMaxIdleTime(time.Second * time.Duration(0))
	db.SetConnMaxLifetime(time.Second * time.Duration(0))
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(0)

	if err != nil {
		fmt.Println(strConn)
		return nil, err
	}

	errPing := db.Ping()

	if errPing != nil {
		fmt.Println("database connection failed, trying new connection.")
		fmt.Println(strConn)

		var i int = 5

		for {

			log := fmt.Sprintf("connection %d of %d", i, 5)
			fmt.Println(log)

			time.Sleep(10 * time.Second)
			i--
			errPing2 := db.Ping()

			if errPing2 == nil {
				return db, nil
			}

			if i == 0 {
				db.Close()
				defer os.Exit(3)
				return nil, err
			}

		}

	}

	return db, err
}
