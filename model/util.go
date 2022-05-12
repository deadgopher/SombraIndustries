package model

import "database/sql"

var eveHQ *sql.DB

func Init(db *sql.DB) {
	eveHQ = db
}
