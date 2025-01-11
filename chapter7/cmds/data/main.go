package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	path := flag.String("path", "./data.sqlite", "path to the database file")
	flag.Parse()

	if _, err := os.Stat(*path); err == nil {
		log.Println("File already exists")
		os.Exit(0)
	}

	db, err := sql.Open("sqlite3", *path+"?_journal=wal")
	if err != nil {
		log.Println("sql.Open", err)
		os.Exit(-1)
	}
	defer func() {
		_ = db.Close()
	}()

	now := time.Now().UTC()
	past := now.Add(-24 * time.Hour).Unix()
	near := now.Add(5 * time.Hour).Unix()
	future := now.Add(5 * 24 * time.Hour).Unix()

	query = fmt.Sprintf(query, past, near, near, future)
	_, err = db.Exec(query)
	if err != nil {
		log.Println("db.Exec", err)
		os.Exit(-1)
	}

	fmt.Println("data created successfully!")
}

var query = `
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    trial_extended BOOLEAN NOT NULL,
    expired_at INTEGER
);

CREATE TABLE payment_methods (
    account_id INTEGER NOT NULL
);

INSERT INTO accounts (id, trial_extended, expired_at) VALUES
(1, false, null),
(2, true, %d),
(3, true, null),
(4, false, %d),
(5, true, %d),
(6, false, %d);

INSERT INTO payment_methods (account_id) VALUES
(1),
(3);
`
