package main

import (
	"bufio"
	"regexp"
	"strings"
	"testing"
)

const text = `2016-09-28 04:30:30, Info CBS    Starting TrustedInstaller initialization.
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
2016-11-09 12:06:05, Info                  CBS    Session: 30554686_2545588473 initialized by client WindowsUpdateAgent.
2016-11-09 12:06:05, Info                  CBS    Read out cached package applicability for package: Package_for_KB2647753~31bf3856ad364e35~amd64~~6.1.4.0, ApplicableState: 112, CurrentState:112
2017-04-09 14:40:09, Info                  CBS    Starting TrustedInstaller initialization.
2017-04-09 14:40:09, Info                  CBS    Loaded Servicing Stack v6.1.7601.23505 with Core: C:\Windows\winsxs\amd64_microsoft-windows-servicingstack_31bf3856ad364e35_6.1.7601.23505_none_681aa442f6fed7f0\cbscore.dll
2017-04-09 14:40:10, Info                  CSI    00000001@2017/4/9:06:40:10.979 WcpInitialize (wcp.dll version 0.0.0.6) called (stack @0x7fef230eb5d @0x7fef25c9b6d @0x7fef259358f @0xff28e97c @0xff28d799 @0xff28db2f)
2017-04-09 14:40:10, Info                  CSI    00000002@2017/4/9:06:40:10.983 WcpInitialize (wcp.dll version 0.0.0.6) called (stack @0x7fef230eb5d @0x7fef2616ade @0x7fef25e2984 @0x7fef2593665 @0xff28e97c @0xff28d799)
2017-04-09 14:40:10, Info                  CSI    00000003@2017/4/9:06:40:10.984 WcpInitialize (wcp.dll version 0.0.0.6) called (stack @0x7fef230eb5d @0x7fef20b8728 @0x7fef20b8856 @0xff28e474 @0xff28d7de @0xff28db2f)
2017-04-09 14:40:10, Info                  CBS    Ending TrustedInstaller initialization.
2017-04-09 14:40:10, Info                  CBS    Starting the TrustedInstaller main loop.
2017-04-09 14:40:10, Info                  CBS    TrustedInstaller service starts successfully.
2017-04-09 14:40:10, Info                  CBS    SQM: Initializing online with Windows opt-in: False
2017-04-09 14:40:10, Info                  CBS    SQM: Cleaning up report files older than 10 days.
2017-04-09 14:40:10, Info                  CBS    SQM: Requesting upload of all unsent reports.
2017-04-09 14:40:10, Info                  CBS    SQM: Failed to start upload with file pattern: C:\Windows\servicing\sqm\*_std.sqm, flags: 0x2 [HRESULT = 0x80004005 - E_FAIL]
2017-04-09 14:40:10, Info                  CBS    SQM: Failed to start standard sample upload. [HRESULT = 0x80004005 - E_FAIL]
2017-04-09 14:40:10, Info                  CBS    SQM: Queued 0 file(s) for upload with pattern: C:\Windows\servicing\sqm\*_all.sqm, flags: 0x6
2017-04-09 14:40:10, Info                  CBS    SQM: Warning: Failed to upload all unsent reports. [HRESULT = 0x80004005 - E_FAIL]
2017-04-09 14:40:10, Info                  CBS    No startup processing required, TrustedInstaller service was not set as autostart, or else a reboot is still pending.
2017-04-09 14:40:10, Info                  CBS    NonStart: Checking to ensure startup processing was not required.
2017-04-09 14:40:10, Info                  CSI    00000004 IAdvancedInstallerAwareStore_ResolvePendingTransactions (call 1) (flags = 00000004, progress = NULL, phase = 0, pdwDisposition = @0x10dfd00
2017-04-09 14:40:11, Info                  CSI    00000005 Creating NT transaction (seq 1), objectname [6]"(null)"
2017-04-09 14:40:11, Info                  CSI    00000006 Created NT transaction (seq 1) result 0x00000000, handle @0x218
2017-04-09 14:40:11, Info                  CSI    00000007@2017/4/9:06:40:11.355 CSI perf trace:
CSIPERF:TXCOMMIT;26151
2017-04-09 14:40:11, Info                  CBS    NonStart: Success, startup processing not required as expected.
2017-04-09 14:40:11, Info                  CBS    Startup processing thread terminated normally
2017-04-09 14:40:11, Info                  CSI    00000008 CSI Store 3116512 (0x00000000002f8de0) initialized
2017-04-09 14:40:11, Info                  CBS    Session: 30585084_589257473 initialized by client lpksetup.
2017-04-09 14:40:17, Info                  CBS    Session: 30585084_589257473 finalized. Reboot required: no [HRESULT = 0x00000000 - S_OK]
2017-04-09 14:40:34, Info                  CBS    Session: 30585084_820720712 initialized by client WinMgmt.
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Warning: Unrecognized packageExtended attribute.
2017-04-09 14:41:50, Info                  CBS    Expecting attribute name [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]
2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]`

func TestLogSplit(t *testing.T) {
	sc := bufio.NewScanner(strings.NewReader(text))

	pat := regexp.MustCompile(defaultDTPattern)
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

func TestCanReadLogEntriesByStartPattern(t *testing.T) {
	entriesOut := readEntries(bufio.NewReader(strings.NewReader(text)), regexp.MustCompile(defaultDTPattern), 100)

	entries := []string{}
	
	for e := range entriesOut {
		entries = append(entries, e)
	}

	if len(entries) != 106 {
		t.Fatalf("expecting 106 rows but was: %d", len(entries))
	}

	if entries[105] != `2017-04-09 14:41:50, Info                  CBS    Failed to get next element [HRESULT = 0x800f080d - CBS_E_MANIFEST_INVALID_ITEM]` {
		t.Fatalf("unexpected: %s", entries[105])
	}
}
