package main

import (
	"bufio"
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type args struct {
	TxSize  int
	InFile  string
	OutFile string
}

func main() {
	var args args
	flag.IntVar(&args.TxSize, "tx", 1000, "specify size of the tx.")

	flag.Parse()

	args.InFile = flag.Arg(0)
	args.OutFile = flag.Arg(1)

	db, err := sql.Open("sqlite3",
		args.OutFile+"?_journal=OFF&_locking=EXCLUSIVE&_sync=OFF")
	if err != nil {
		log.Fatalf("failed to open db: %q", err)
	}

	defer db.Close()

	createTable := `
		create virtual table if not exists logs using fts5(
			time UNINDEXED,
			entry,
			prefix=3,
			tokenize = "unicode61 tokenchars '-_.'"
		);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("failed to create table: %q", err)
	}

	infile, err := os.Open(args.InFile)
	if err != nil {
		log.Fatalf("failed to open input file: %q", err)
	}

	scanner := bufio.NewScanner(infile)
	// scanner.Buffer(make([]byte, 0, 1000*1000*1000), 1000*1000*1000)
	cnt := 0
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin: %q", err)
	}

	stmt, err := tx.Prepare(`insert into logs values (?, ?)`)
	if err != nil {
		log.Fatalf("failed to prepare stmt: %q", err)
	}

	for scanner.Scan() {
		cnt++
		entry := scanner.Text()
		if len(entry) < 19 {
			continue
		}

		_, err = stmt.Exec(entry[:19], entry)
		if err != nil {
			log.Fatalf("failed to insert: %q", err)
		}
		if cnt%args.TxSize == 0 {
			tx.Commit()

			tx, err = db.Begin()
			if err != nil {
				log.Fatalf("failed to begin: %q", err)
			}

			stmt, err = tx.Prepare(`insert into logs values (?, ?)`)
			if err != nil {
				log.Fatalf("failed to prepare stmt: %q", err)
			}
		}
	}

	tx.Commit()
}
