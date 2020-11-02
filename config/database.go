package database

import (
	"database/sql"

	_ "github.com/jinzhu/gorm/dialects/mssql"    // MSSql Driver
	_ "github.com/jinzhu/gorm/dialects/mysql"    // MySql Driver
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres Driver

	"fmt"
	"strings"
)

type DBConn interface {
	Connect() (*sql.DB, error)
}

type DBCred struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Connect to database
func (dbCred DBCred) Connect() (*sql.DB, error) {
	var db *sql.DB
	var err error

	connectionStrings := []string{
		fmt.Sprintf("mysql~%s", fmt.Sprintf(
			"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			dbCred.DBUser,
			dbCred.DBPassword,
			dbCred.DBHost,
			dbCred.DBPort,
			dbCred.DBName,
		)),
		fmt.Sprintf("postgres~%s", fmt.Sprintf(
			"host=%s port=%s user=%s sslmode=disable dbname=%s password=%s",
			dbCred.DBHost,
			dbCred.DBPort,
			dbCred.DBUser,
			dbCred.DBName,
			dbCred.DBPassword,
		)),
		fmt.Sprintf("mssql~%s", fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s",
			dbCred.DBUser,
			dbCred.DBPassword,
			dbCred.DBHost,
			dbCred.DBPort,
			dbCred.DBName,
		)),
	}

	for index := range connectionStrings {
		drv := strings.Split(connectionStrings[index], "~")
		if dbCred.DBDriver == drv[0] {
			db, err = sql.Open(dbCred.DBDriver, drv[1])
			if err != nil {
				return &sql.DB{}, err
			}
			break
		}
	}

	fmt.Println("Database Connected")
	
	return db, nil
}