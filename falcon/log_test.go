package falcon

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestLog (t *testing.T) {
	onCallGroupName := "testPoint"
	LogPath =  "../logs/" + onCallGroupName + "_" + time.Now().Format("20060102150405") + ".txt"

	logFile, err := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("记录日志失败:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(0)
	//log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println("这是一条很普通的日志。")
	log.SetPrefix("[小王子]")
	log.Println("这是一条很普通的日志。")
}
