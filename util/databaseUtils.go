package util

import (
	"database/sql"
	"sync"
	"time"
)

var (
	LOCAL    bool    = true
	Database *sql.DB = nil
	mutex    sync.Mutex
	MyTx     *sql.Tx
)

func GetDatabase() *sql.DB {
	return Database
}

func StartDatabase() {
	var db *sql.DB
	var err error

	mutex.Lock()

	if LOCAL {
		// LOCAL
		db, err = sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")
	} else {
		// PROD
		db, err = sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	}

	if err != nil {
		panic(err)
	}
	errPing := db.Ping()
	if errPing != nil {
		panic("Erro ao acessar o banco de dados!!!")
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Duration(600) * time.Second)

	Database = db
	mutex.Unlock()
}
