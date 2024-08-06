package falcon

import (
	"fmt"
	"log"
	"os"
)

func (f *FALCON) Logging (prefix string, content *[]byte) bool {
	logFile, err := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("记录日志失败:", err)
		return false
	}
	log.SetOutput(logFile)
	log.SetFlags(0)
	//log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	//log.Println("这是一条很普通的日志。")
	log.SetPrefix(prefix + "++")
	log.Println(string(*content))

	return true
}
