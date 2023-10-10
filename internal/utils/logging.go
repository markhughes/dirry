package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/markhughes/dirry/internal/consts"
)

const gray = "\033[37m"
const white = "\033[97m"
const reset = "\033[0m"

var EnabledDebugCategoriesMap map[string]bool = make(map[string]bool)
var EnabledDebugAll bool = false

var EnableLogging bool = false

// log ifle name sohudl be dirry_<timestamp>.log
var LogPrefix string = "dirry_" + time.Now().Format("2006-01-02-15.04.05")
var LogfileName string = LogPrefix + ".log"
var LogfileExtendDirectory string = path.Join(LogPrefix, "cat")

func SaveLog(logType string, category string, format string, a ...interface{}) {
	if !EnableLogging {
		return
	}

	// savees all logs to logs/<LogFileName>
	// saves "types" into logs/<LogfileExtendDirectory>/<logType>/<category>.log

	_, fullFilePath1, line1, _ := runtime.Caller(3)

	_, fullFilePath2, line2, _ := runtime.Caller(2)

	var logLine = fmt.Sprintf("[%s] %s - [%s:%d]>[%s:%d]: ", logType, category, path.Base(fullFilePath1), line1, path.Base(fullFilePath2), line2)
	logLine += fmt.Sprintf(format, a...)

	// Define the paths for the log files.
	logFilePath := filepath.Join(consts.LogsDir, LogPrefix, LogfileName)
	logTypePath := filepath.Join(consts.LogsDir, LogfileExtendDirectory, strings.TrimSpace(logType), strings.TrimSpace(category)+".log")

	for _, filePath := range []string{logFilePath, logTypePath} {
		// Ensure the directory for the log file exists.
		dirPath := filepath.Dir(filePath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			fmt.Println("Failed to create directory:", err)
			return
		}

		// Create (or open) the log file.
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Failed to open log file:", err)
			return
		}
		defer file.Close()

		// Write the log line to the file.
		if _, err := file.WriteString(logLine + "\n"); err != nil {
			fmt.Println("Failed to write to log file:", err)
			return
		}
	}
}

func InfoMsg(category string, format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	SaveLog("INFO", category, format, a...)
}

func WarnMsg(category string, format string, a ...interface{}) {
	fmt.Printf("\033[33m⚠️\033[0m "+format+"\n", a...)
	SaveLog("WARN", category, format, a...)
}

func SuccessMsg(category string, format string, a ...interface{}) {
	fmt.Printf("\033[32m✔\033[0m "+format+"\n", a...)
	SaveLog("INFO", category, format, a...)
}

func ErrorMsg(category string, format string, a ...interface{}) {
	fmt.Printf("\033[31m✖\033[0m "+format+"\n", a...)
	SaveLog("ERRO", category, format, a...)
}

func DebugMsg(category string, format string, a ...interface{}) {
	SaveLog("DBUG", category, format, a...)
	if !EnabledDebugAll {
		_, found := EnabledDebugCategoriesMap[category]
		if !found {
			return
		}
	}

	_, fullFilePath1, line1, _ := runtime.Caller(2)
	fileName1 := fullFilePath1[strings.LastIndex(fullFilePath1, "/")+1:]

	paddedLineNumber1 := fmt.Sprintf("%04d", line1)

	nonZeroIndex1 := strings.IndexFunc(paddedLineNumber1, func(r rune) bool { return r != '0' })

	if nonZeroIndex1 == -1 {
		nonZeroIndex1 = len(paddedLineNumber1)
	}

	coloredPaddedLineNumber1 := gray + paddedLineNumber1[:nonZeroIndex1] + reset
	if nonZeroIndex1 < len(paddedLineNumber1) {
		coloredPaddedLineNumber1 += white + paddedLineNumber1[nonZeroIndex1:] + reset
	}

	//
	_, fullFilePath2, line2, _ := runtime.Caller(1)
	fileName2 := fullFilePath2[strings.LastIndex(fullFilePath2, "/")+1:]

	paddedLineNumber2 := fmt.Sprintf("%04d", line2)

	nonZeroIndex2 := strings.IndexFunc(paddedLineNumber2, func(r rune) bool { return r != '0' })

	if nonZeroIndex2 == -1 {
		nonZeroIndex2 = len(paddedLineNumber2)
	}

	coloredPaddedLineNumber2 := gray + paddedLineNumber2[:nonZeroIndex2] + reset
	if nonZeroIndex2 < len(paddedLineNumber2) {
		coloredPaddedLineNumber2 += white + paddedLineNumber2[nonZeroIndex2:] + reset
	}

	fmt.Printf("%s - [%s:%s]>[%s:%s]: ", category, fileName1, coloredPaddedLineNumber1, fileName2, coloredPaddedLineNumber2)
	fmt.Printf(format+reset+"\n", a...)
}
