package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (    
	id INTEGER PRIMARY KEY AUTOINCREMENT,    
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX date ON scheduler (date);
`

var db *sql.DB

func Init(dbFile string) error {
	var install bool
	_, err := os.Stat(dbFile)

	if err != nil {
		install = true
	}
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install == true {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
