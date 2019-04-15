package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/revel/revel"
	"log"
)

var Handle *sqlx.DB

func Initialize() {
	driver := revel.Config.StringDefault("db.driver", "mysql")
	connectString := revel.Config.StringDefault("db.connect",
		"root:@(localhost:3306)/blog")

	db, err := sqlx.Connect(driver, connectString)
	if err != nil {
		log.Fatal("[!] DB Err: ", err)
	}

	Handle = db
}