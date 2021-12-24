package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/HaoxuanXu/TradingBot/db"
)

func SetLogging(assetType string) *os.File {
	dt := time.Now()
	logName := fmt.Sprintf("%d-%d-%d", dt.Year(), dt.Month(), dt.Day()) + "_" + "TradingLog.log"
	fullLogPath := db.MapLogPath(assetType) + logName
	logFile, err := os.OpenFile(fullLogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	multiWrite := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiWrite)
	log.Printf("logging the trading record to %s\n", fullLogPath)
	return logFile
}
