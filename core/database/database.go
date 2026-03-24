package database

import (
	"database/sql"
	"gvm/core/config"
	"log"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"
)

func GetDbVersion() string {
	var version string
	var dbFilePath = config.ConfigDirectory + string(os.PathSeparator) + "gvm.db"

	var db, err = sql.Open("sqlite3", "file:"+dbFilePath)

	if err != nil {
		log.Printf("Error opening database: %v", err)
		return ""
	}
	defer db.Close()

	db.QueryRow(`SELECT sqlite_version()`).Scan(&version)
	return version
}
