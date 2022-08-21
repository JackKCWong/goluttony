package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

type args struct {
	TxSize          int
	DatetimePattern string
	InFile          string
	OutFile         string
}

func main() {
	var args args
	flag.IntVar(&args.TxSize, "tx", 1000, "specify size of the tx.")
	flag.StringVar(&args.DatetimePattern, "dt", `\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}:\d{2}`, "specify datetime regex pattern to use. The pattern is used to split log entries.")

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
			datetime,
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
	pat, err := regexp.Compile(args.DatetimePattern)
	if err != nil {
		log.Fatalf("invalid datatime format: %q", err)
	}
	scanner.Split(logEntrySplitter(pat))
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
			log.Printf("skip short line: %s", entry)
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

func logEntrySplitter(dtPattern *regexp.Regexp) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		loc := dtPattern.FindIndex(data)
		if loc == nil {
			if atEOF {
				return 0, nil, fmt.Errorf("unexpected EOF: buffer=%s", data)
			}

			return 0, nil, nil
		}

		// assumes data always start with a dtPattern
		if loc[0] == 0 {
			if atEOF {
				// last line
				return len(data), bytes.TrimRight(data, "\n"), nil
			}

			// first line
			loc = dtPattern.FindIndex(data[1:]) // find next entry start
			if loc == nil {
				return 0, nil, nil
			}

			// fix offset by 1 above
			loc[0] += 1
			loc[1] += 1

			if data[loc[0]-1] == '\n' {
				// a good match
				return loc[0],
					data[0 : loc[0]-1], // don't return \n at the end
					nil
			} else {
				return 0, nil, nil
			}
		} else {
			// read more
			return 0, nil, fmt.Errorf("unexpected start of entry: %s", data) 
		}
	}
}
