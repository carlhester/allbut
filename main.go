package main

import (
	"log"
	"os"

	"github.com/carlhester/allbut/allbut"
)

func main() {
	app, err := allbut.Setup(os.Args[1:])
	if err != nil {
		log.Println(err)
		allbut.PrintUsageAndExit()
	}
	err = app.Run()
	if err != nil {
		log.Println(err)
		allbut.PrintUsageAndExit()
	}

}
