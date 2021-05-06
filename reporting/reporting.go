package reporting

import (
	"fmt"
	"os"
	"strings"
	"time"

	humbug "github.com/bugout-dev/humbug/go/pkg"

	"github.com/hashicorp/nomad/version"
)

var SessionID string = fmt.Sprintf("%d", time.Now().Unix())
var ClientID string = os.Getenv("NOMAD_REPORTS_EMAIL")

const ReporterToken string = "6924ac70-f2b4-4444-a26b-7db8c15bf694"

var UserConsentState bool = false

func CheckUserConsent() bool {
	return UserConsentState
}

var ReportingConsent *humbug.HumbugConsent = humbug.CreateHumbugConsent(CheckUserConsent)

var Reporter *humbug.HumbugReporter

func InitializeReporter() *humbug.HumbugReporter {
	if Reporter != nil {
		return Reporter
	}

	var err error
	Reporter, err = humbug.CreateHumbugReporter(ReportingConsent, ClientID, SessionID, ReporterToken)
	if err != nil {
		panic(err)
	}

	nomadVersion := version.GetVersion().VersionNumber()
	Reporter.Tag("version", nomadVersion)

	return Reporter
}

func CommandInvokedReport(args []string) humbug.Report {
	report := humbug.SystemReport()
	command := strings.Join(args, " ")
	report.Title = fmt.Sprintf("Command invoked: %s", command)
	report.Content = fmt.Sprintf("Command: `%s`\nInvoked at: `%v`\n", command, time.Now())
	return report
}

func CommandCompletedReport(args []string, code int) humbug.Report {
	report := humbug.SystemReport()
	command := strings.Join(args, " ")
	report.Title = fmt.Sprintf("Command exited (%d): %s", code, command)
	report.Content = fmt.Sprintf("Command: `%s`\nExited at: `%v`\nExit code: `%d`\n", command, time.Now(), code)
	report.Tags = append(report.Tags, fmt.Sprintf("exit:%d", code))
	return report
}
