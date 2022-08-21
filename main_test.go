package main

import (
	"bufio"
	"regexp"
	"strings"
	"testing"
)

func TestLogSplit(t *testing.T) {
	sc := bufio.NewScanner(strings.NewReader(`2016-09-28 04:30:30, Info CBS    Starting TrustedInstaller initialization.
2016-09-28 04:30:30, Info CBS    Loaded Servicing Stack v6.1.7601.23505 with Core: 
C:\Windows\winsxs\amd64_microsoft-windows-servicingstack_31bf3856ad364e35_6.1.7601.23505_none_681aa442f6fed7f0\cbscore.dll
2016-09-28 04:30:31, Info CSI    00000001@2016/9/27:20:30:31.455 WcpInitialize (wcp.dll version 0.0.0.6) called (stack @0x7fed806eb5d @0x7fef9fb9b6d @0x7fef9f8358f @0xff83e97c @0xff83d799 @0xff83db2f)
2016-09-28 04:30:31, Info CSI    00000002@2016/9/27:20:30:31.458 WcpInitialize (wcp.dll version 0.0.0.6) called (stack @0x7fed806eb5d @0x7fefa006ade @0x7fef9fd2984 @0x7fef9f83665 @0xff83e97c @0xff83d799)
2016-09-28 04:30:31, Info CSI    00000003@2016/9/27:20:30:31.458 WcpInitialize (wcp.dll version 0.0.0.6) called (stack @0x7fed806eb5d @0x7fefa1c8728 @0x7fefa1c8856 @0xff83e474 @0xff83d7de @0xff83db2f)
2016-09-28 04:30:31, Info CBS    Ending TrustedInstaller initialization.
Info CBS    Starting the TrustedInstaller main loop.
Info CBS    TrustedInstaller service starts successfully.
Info CBS    SQM: Initializing online with Windows opt-in: False
Info CBS    SQM: Cleaning up report files older than 10 days.
`))

	pat := regexp.MustCompile(`\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}:\d{2}`)
	sc.Split(logEntrySplitter(pat))

	var entries []string

	for sc.Scan() {
		entries = append(entries, sc.Text())
	}

	if len(entries) != 6 {
		t.Fatalf("expecting 6 entries but was %d", len(entries))
	}

	if entries[0] != `2016-09-28 04:30:30, Info CBS    Starting TrustedInstaller initialization.` {
		t.Fatalf("was %s", entries[0])
	}

	if entries[5] != `2016-09-28 04:30:31, Info CBS    Ending TrustedInstaller initialization.
Info CBS    Starting the TrustedInstaller main loop.
Info CBS    TrustedInstaller service starts successfully.
Info CBS    SQM: Initializing online with Windows opt-in: False
Info CBS    SQM: Cleaning up report files older than 10 days.` {
		t.Fatalf("was %s", entries[0])
	}
}
