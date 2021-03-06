package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/HaoxuanXu/TradingBot/db"
)

func SetLogging(assetType string) *os.File {
	dt := time.Now()
	logName := fmt.Sprintf("%d-%d-%d_%s_TradingLog.log", dt.Year(), dt.Month(), dt.Day(), assetType)
	fullLogPath := db.MapLogPath(assetType) + logName
	monitorLogPath := db.MapLogPath("monitor") + strings.Title(assetType) + "TradingLog.log"
	logFile, err := os.OpenFile(fullLogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	monitorLog, err := os.OpenFile(monitorLogPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	multiWrite := io.MultiWriter(logFile, monitorLog)
	log.SetOutput(multiWrite)
	log.Printf("logging the trading record to %s\n", fullLogPath)
	return logFile
}
