package main

import (
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/zhengxiaoyao0716/oss2bdc/oss2bdc"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	defer oss2bdc.ShutDown()

	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
		log.Println(dir)
	} else {
		dir = time.Unix(time.Now().Unix()-24*60*60, 0).Format("2006/01-02")
	}

	log.Println("Running: ", dir)
	oss2bdc.Download(dir)
	log.Println("Download completed.")
	oss2bdc.Compress(dir)
	log.Println("Compress completed.")
}
