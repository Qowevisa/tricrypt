package main

import (
	"log"
	"os"

	"git.qowevisa.me/Qowevisa/gotell/tui"
)

func main() {
	logFile, err := os.OpenFile("ctui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Printf("Start")

	err = tui.UI.Run()
	if err != nil {
		panic(err)
	}
}
