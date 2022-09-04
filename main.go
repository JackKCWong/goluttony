package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/profile"
)

type args struct {
	TxSize          int
	BufSize         int
	DatetimePattern string
	TokenChars      string
	Profile         string
	InFile          string
	DbURL           string
}

const defaultDTPattern = `2006-01-02 15:04:05`

type Line struct {
	Time time.Time
	Raw  string
}

func main() {
	var args args
	flag.IntVar(&args.TxSize, "tx", 1000, "specify size of the tx.")
	flag.IntVar(&args.BufSize, "sz", 64, "specify max size of a line in kb.")
	flag.StringVar(&args.DatetimePattern, "dt", defaultDTPattern, "specify datetime layout to use.")
	flag.StringVar(&args.TokenChars, "tok", "_", "extra token chars for fts")
	flag.StringVar(&args.Profile, "prof", "", "enable either cpu|mem profiling")

	flag.Parse()

	switch args.Profile {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	case "":
		// do nothing
	default:
		log.Fatalf("-prof must be one of cpu|mem, but was %s.", args.Profile)
	}

	args.InFile = flag.Arg(0)
	args.DbURL = flag.Arg(1)

	db, err := sql.Open("sqlite3", args.DbURL)
	if err != nil {
		log.Fatalf("failed to open db: %q", err)
	}

	defer db.Close()

	createTable := `
		create table if not exists logs (
			time integer NOT NULL,
			raw text NOT NULL
		);
		create virtual table if not exists logs_fts using fts5(
			raw,
			tokenize="unicode61 tokenchars '%s'",
			content="logs"
		);
		create trigger logs_fts_insert after insert on logs 
		begin
			insert into logs_fts (rowid, raw) values (new.rowid, new.raw);
		end;
	`

	_, err = db.Exec(fmt.Sprintf(createTable, args.TokenChars))
	if err != nil {
		log.Fatalf("failed to create table: %q", err)
	}

	infile, err := os.Open(args.InFile)
	if err != nil {
		log.Fatalf("failed to open input file: %q", err)
	}

	br := bufio.NewReaderSize(infile, args.BufSize*1024)
	r, _, err := br.ReadRune()
	if err != nil {
		log.Fatalf("failed to read file: %q", err)
	}

	// utf8 BOM check
	if r != '\uFEFF' {
		br.UnreadRune()
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin: %q", err)
	}

	stmt, err := db.Prepare(`insert into logs values (?, ?)`)
	if err != nil {
		log.Fatalf("failed to prepare stmt: %q", err)
	}
	defer stmt.Close()

	txstmt := tx.Stmt(stmt)
	start := time.Now()
	cnt := 0

	for entry := range readEntries(br, args.DatetimePattern, args.TxSize+1) {
		cnt++
		_, err = txstmt.Exec(entry.Time, entry.Raw)
		if err != nil {
			log.Fatalf("failed to insert: %q", err)
		}
		if cnt%args.TxSize == 0 {
			// the order of Commit & Close doesn't seem to matter
			err = tx.Commit()
			if err != nil {
				log.Fatalf("failed to commit tx: %q", err)
			}

			err = txstmt.Close()
			if err != nil {
				log.Fatalf("failed to close stmt: %q", err)
			}

			tx, err = db.Begin()
			if err != nil {
				log.Fatalf("failed to begin tx: %q", err)
			}

			txstmt = tx.Stmt(stmt)
		}
	}

	tx.Commit()
	end := time.Now()
	log.Printf("%d rows committed in %s", cnt, end.Sub(start))
	log.Printf("%d rows/s", cnt/int(end.Sub(start).Seconds()))
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
		sb.WriteByte('\n')
		sb.Write(line)
	}

	return sb.String(), nil
}

func readEntries(br *bufio.Reader, dtLayout string, bufsize int) chan Line {
	out := make(chan Line, bufsize)
	go func() {
		defer close(out)
		var entry strings.Builder
		var line1 string
		var lineN string
		var ioErr error

		line1, ioErr = readFullLine(br)
		for {
			if len(line1) > 0 {
				dt, parseErr := parseTimestamp(dtLayout, line1)
				if parseErr != nil {
					log.Fatalf("unexpected starting of line: %s", line1)
				}

				entry.Reset()
				entry.WriteString(line1)

				for {
					lineN, ioErr = readFullLine(br)
					if len(lineN) > 0 {
						_, parseErr = parseTimestamp(dtLayout, lineN)
						if parseErr != nil {
							entry.WriteByte('\n')
							entry.WriteString(lineN)
						} else {
							line1 = lineN
							break
						}
					} else {
						if ioErr == io.EOF {
							break
						}

						if ioErr != nil {
							log.Fatalf("failed to read next line: %q", ioErr)
						}
					}
				}

				out <- Line{
					Time: dt,
					Raw:  entry.String(),
				}
			}

			if ioErr == io.EOF {
				break
			}

			if ioErr != nil {
				log.Fatalf("failed to read line: %q", ioErr)
			}
		}
	}()

	return out
}

func parseTimestamp(layout, line string) (time.Time, error) {
	if len(line) > len(layout) {
		return time.Parse(layout, line[0:len(layout)])
	}

	return time.Time{}, fmt.Errorf("input shorter than layout: %s", line)
}
