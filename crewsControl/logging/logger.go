package logging

// NOTE at the time of writing log/syslog is not implemented for windows
import (
	"log"
	"log/syslog"
)

var SysLogger *syslog.Writer

// Redirects the output of log.* functions to write to syslog with the tag LOG_INFO
func InitLogger() {
	SysLogger, _ = syslog.New(syslog.LOG_INFO, "cot")
	log.SetOutput(SysLogger)
}

func TestLog() {
	log.Println("Test log entry")
	log.Fatalf("Test fatal")
}
