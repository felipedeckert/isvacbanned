package util

import (
	"database/sql"
	"os"
	"sync"
	"time"
)

var (
	LOCAL    bool    = true
	Database *sql.DB = nil
	mutex    sync.Mutex
	MyTx     *sql.Tx
)

//GetDatabase returns a *sql.DB, that should have been initialized before by StartDatabase()
func GetDatabase() *sql.DB {
	return Database
}

//StartDatabase opens a database connection
func StartDatabase() {
	var (
		db  *sql.DB
		err error
	)

	credentials := os.Getenv("DATABASE_CREDENTIAL")

	mutex.Lock()

	if LOCAL {
		// LOCAL
		db, err = sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")
	} else {
		// PROD
		db, err = sql.Open("mysql", credentials+"@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	}

	if err != nil {
		panic(err)
	}
	errPing := db.Ping()
	if errPing != nil {
		panic("Erro ao acessar o banco de dados!!!")
	}

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Duration(600) * time.Second)

	Database = db
	mutex.Unlock()
}
