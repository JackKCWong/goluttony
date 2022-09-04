package main

import (
	"bufio"
	"strings"
	"testing"
)

const text = `2016-09-28 04:30:30, Info CBS    Starting TrustedInstaller initialization.
2016-09-28 04:30:30, Info CBS    Loaded Servicing Stack v6.1.7601.23505 with Core: 
C:\Windows\winsxs\amd64_microsoft-windows-servicingstack_31bf3856ad364e35_6.1.7601.23505_none_681aa442f6fed7f0\cbscore.dll
2016-09-28 04:30:31, Info CSI    00000003@2016/9/27:20:30:31.458 WcpInitialize (wcp.dll version 0.0.0.6) called (stack @0x7fed806eb5d @0x7fefa1c8728 @0x7fefa1c8856 @0xff83e474 @0xff83d7de @0xff83db2f)
2016-09-28 04:30:32, Info CBS    Ending TrustedInstaller initialization.
Info CBS    Starting the TrustedInstaller main loop.
Info CBS    TrustedInstaller service starts successfully.
Info CBS
Info CBS    SQM: Cleaning up report files older than 10 days.
2017-04-09 14:40:10, Info                  CBS    NonStart: Checking to ensure startup processing was not required.
2017-04-09 14:40:11, Info                  CSI    00000007@2017/4/9:06:40:11.355 CSI perf trace:
CSIPERF:TXCOMMIT;26151
`

func TestCanReadLogEntriesByStartPattern(t *testing.T) {
	entriesOut := readEntries(bufio.NewReader(strings.NewReader(text)), defaultDTPattern, 100)

	entries := []Line{}

	for e := range entriesOut {
		entries = append(entries, e)
	}

	if len(entries) != 6 {
		t.Fatalf("expecting 8 rows but was: %d", len(entries))
	}

	if !strings.HasSuffix(entries[3].Raw, `Cleaning up report files older than 10 days.`) {
		t.Fatalf("unexpected: %s", entries[3].Raw)
	}
}
