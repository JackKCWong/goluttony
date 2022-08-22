package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/profile"
)

type args struct {
	TxSize          int
	DatetimePattern string
	Profile         bool
	InFile          string
	OutFile         string
}

const defaultDTPattern = `\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}:\d{2}`

func main() {
	var args args
	flag.IntVar(&args.TxSize, "tx", 1000, "specify size of the tx.")
	flag.StringVar(&args.DatetimePattern, "dt", defaultDTPattern, "specify datetime regex pattern to use. The pattern is used to split log entries.")
	flag.BoolVar(&args.Profile, "prof", false, "enable profiling")

	flag.Parse()

	if args.Profile {
		defer profile.Start(
			profile.MemProfile,
			profile.CPUProfile,
			profile.ProfilePath(".")).
			Stop()
	}

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
			entry,
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

	br := bufio.NewReader(infile)
	r, _, err := br.ReadRune()
	if err != nil {
		log.Fatalf("failed to read file: %q", err)
	}

	// utf8 BOM check
	if r != '\uFEFF' {
		br.UnreadRune()
	}

	pat, err := regexp.Compile(args.DatetimePattern)
	if err != nil {
		log.Fatalf("invalid datatime format: %q", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin: %q", err)
	}

	stmt, err := tx.Prepare(`insert into logs values (?)`)
	if err != nil {
		log.Fatalf("failed to prepare stmt: %q", err)
	}

	var cnt = 0
	for entry := range readEntries(br, pat, args.TxSize+1) {
		_, err = stmt.Exec(entry)
		if err != nil {
			log.Fatalf("failed to insert: %q", err)
		}
		if cnt%args.TxSize == 0 {
			tx.Commit()

			tx, err = db.Begin()
			if err != nil {
				log.Fatalf("failed to begin: %q", err)
			}

			stmt, err = tx.Prepare(`insert into logs values (?)`)
			if err != nil {
				log.Fatalf("failed to prepare stmt: %q", err)
			}
		}
	}

	tx.Commit()
}

func readFullLine(in *bufio.Reader) (string, error) {
	var sb *strings.Builder
	line, prefix, err := in.ReadLine()

	if err != nil {
		return "", err
	}

	if !prefix {
		return string(line), nil
	}

	sb = new(strings.Builder)
	sb.Write(line)
	for prefix {
		line, prefix, err = in.ReadLine()
		if err != nil {
			return "", err
		}
		sb.Write(line)
	}

	return sb.String(), nil
}

func readEntries(br *bufio.Reader, pat *regexp.Regexp, bufsize int) chan string {
	out := make(chan string, bufsize)
	go func() {
		defer close(out)
		var entry strings.Builder

		line, err := readFullLine(br)
		if line == "" {
			if err != io.EOF {
				log.Fatal(err)
			}
		}

		entry.WriteString(line)

	outer:
		for {
			loc := pat.FindStringIndex(entry.String())
			if loc == nil || loc[0] != 0 {
				log.Fatalf("unexpected starting line: %s", entry.String())
			}

			for {
				line, err = readFullLine(br)
				if err == io.EOF {
					out <- entry.String()
					break outer
				}

				if err != nil {
					log.Fatal(err)
				}

				loc = pat.FindStringIndex(line)
				if loc == nil || loc[0] != 0 {
					entry.WriteByte('\n')
					entry.WriteString(line)
				} else {
					// a new entry start
					out <- entry.String()
					entry.Reset()
					entry.WriteString(line)
					break
				}
			}
		}

	}()

	return out
}

func logEntrySplitter(dtPattern *regexp.Regexp) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		loc := dtPattern.FindIndex(data)
		if loc == nil {
			if atEOF {
				return 0, nil, fmt.Errorf("unexpected EOF: buffer=%s", data)
			}

			log.Printf("read more 1: %d", len(data))
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
				log.Printf("read more 2: %d", len(data))
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
				log.Printf("read more 3: %d", len(data))
				return 0, nil, nil
			}
		} else {
			// read more
			return 0, nil, fmt.Errorf("unexpected start of entry: %s,%v", data, loc)
		}
	}
}
