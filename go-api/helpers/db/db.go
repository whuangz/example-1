package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	config "github.com/whuangz/go-example/go-api/config"
)

func init() {
	if config.IS_DEBUG_MODE {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func Configure(dataSource string) *sqlx.DB {
	if dataSource == "" {
		dataSource = config.URI
	}
	client := connect(dataSource)
	return client
}

func connect(dataSource string) *sqlx.DB {
	conn, err := sqlx.Connect("mysql", dataSource)

	if err != nil {
		logrus.Error(err)
	} else {
		return conn
	}
	return nil
}

// Transaction is an interface that models the standard transaction in
// `database/sql`.
//
// To ensure `TxFn` funcs cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(Transaction) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(db *sqlx.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}

		db.Close()
	}()

	err = fn(tx)
	return err
}
