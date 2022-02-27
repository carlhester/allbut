package main

import (
	"log"
	"os"

	"github.com/carlhester/allbut/allbut"
)

func main() {
	app := allbut.New()
	ab, err := app.Setup(os.Args[1:])
	if err != nil {
		log.Println(err)
		allbut.PrintUsageAndExit()
	}

	err = ab.Run()
	if err != nil {
		log.Println(err)
		allbut.PrintUsageAndExit()
	}

}
