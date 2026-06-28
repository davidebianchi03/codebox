package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gitlab.com/codebox4073715/codebox/config"
)

const logsFilePrefix = "codebox.logs.json."

/*
retrieve the path of the folder where log files are stored
*/
func getLogsDir() string {
	dirpath := filepath.Join(
		config.Environment.UploadsPath,
		"system-logs",
	)

	os.MkdirAll(dirpath, 0700)

	return dirpath
}

/*
get the path of the current log file
*/
func getLogFile() string {
	dirpath := getLogsDir()

	return filepath.Join(
		dirpath,
		fmt.Sprintf("%s%s", logsFilePrefix, time.Now().Format(time.DateOnly)),
	)
}

/*
format log to json string
*/
func formatLog(
	module string,
	function string,
	level string,
	msg string,
) string {
	data, _ := json.Marshal(map[string]string{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"module":    module,
		"function":  function,
		"level":     level,
		"log":       msg,
	})
	return string(data)
}

/*
log a message with 'info' level
*/
func Info(msg string, args ...any) {
	logFile, err := os.OpenFile(getLogFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

	logFile.Write(
		[]byte(
			formatLog(
				funcName[:lastDot],
				funcName[lastDot+1:],
				"info",
				fmt.Sprintf(msg, args...),
			),
		),
	)

	logFile.Close()
}

/*
log a message with 'error' level
*/
func Error(msg string, args ...any) {
	logFile, err := os.OpenFile(getLogFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

	logFile.Write(
		[]byte(
			formatLog(
				funcName[:lastDot],
				funcName[lastDot+1:],
				"error",
				fmt.Sprintf(msg, args...),
			),
		),
	)

	logFile.Close()
}

/*
log a message with 'debug' level
*/
func Debug(msg string, args ...any) {
	logFile, err := os.OpenFile(getLogFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

	logFile.Write(
		[]byte(
			formatLog(
				funcName[:lastDot],
				funcName[lastDot+1:],
				"debug",
				fmt.Sprintf(msg, args...),
			),
		),
	)

	logFile.Close()
}

/*
log a message with 'warn' level
*/
func Warn(msg string, args ...any) {
	logFile, err := os.OpenFile(getLogFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

	logFile.Write(
		[]byte(
			formatLog(
				funcName[:lastDot],
				funcName[lastDot+1:],
				"warn",
				fmt.Sprintf(msg, args...),
			),
		),
	)

	logFile.Close()
}

/*
rotates logs, remove file with logs older than 7 days
*/
func RotateLogs() error {
	dirpath := getLogsDir()
	entries, _ := os.ReadDir(dirpath)
	for _, entry := range entries {
		// parse filename
		if !entry.IsDir() &&
			strings.Index(entry.Name(), logsFilePrefix) == 0 {
			dateStr := strings.TrimPrefix(entry.Name(), logsFilePrefix)
			date, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				return err
			}

			// remove logs older than 7 days
			if time.Since(date) > 7*24*time.Hour {
				os.RemoveAll(filepath.Join(dirpath, entry.Name()))
			}
		}
	}

	return nil
}
