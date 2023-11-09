package logging

// NOTE at the time of writing log/syslog is not implemented for windows
import (
	"log"
	"log/syslog"
	"os"
	"strconv"
	"time"
)

var SysLogger *syslog.Writer

// Redirects the output of log.* functions to write to syslog with the tag LOG_INFO
func InitLogger() {
	SysLogger, _ = syslog.New(syslog.LOG_INFO, "cot")
	log.SetOutput(SysLogger)
	//  /var/run/cot.pid
	myPid := os.Getpid()
	pidStr := strconv.Itoa(myPid)
	pidBytes := []byte(pidStr)
	err := os.WriteFile("/var/run/cot.pid", pidBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

// Function that will append a string to the designated file
func LogRequest(request string) {
	// /home/crew/crew/crewsOutTherecrewsoutthere/requestLog.txt
	myFile, err := os.OpenFile("/home/crew/crew/crewsoutthere/requestLog.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Could not open file: %s\n", err)
	}

	defer myFile.Close()
	curTime := time.Now().UTC().String()
	writeString := curTime + "    " + request + "\n"
	_, err = myFile.WriteString(writeString)
	if err != nil {
		log.Fatalf("Could not write to request log file: %s\n", err)
	}
}

func TestLog() {
	log.Println("Test log entry")
	log.Fatalf("Test fatal")
}
