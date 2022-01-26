package main

import (
	"fmt"
	"io/fs"
	"os"
)

func main() {
	saveFile, force := parseArgs(os.Args)
	fmt.Println(saveFile, force)

}

func parseArgs(args []string) (string, bool) {
	// Get the file that we're protecting
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("usage: allbut [-f] filename")
	}
	force := false
	saveFile := ""
	for _, arg := range args {
		if arg == "-f" {
			force = true
		} else {
			saveFile = arg
		}
	}
	return saveFile, force

	// // Enumerate fileList
	// fileList, err := ioutil.ReadDir("./")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // first make sure the target file exists
	// if !saveFileExists(saveFile, fileList) {
	// 	fmt.Printf("fatal: file [%s] not found in current directory\n", saveFile)
	// 	os.Exit(0)
	// }

	// // Iterate and decide
	// for _, file := range fileList {
	// 	if file.Name() != saveFile {
	// 		if *force {
	// 			e := os.Remove(file.Name())
	// 			if e != nil {
	// 				log.Fatal(e)
	// 			}
	// 		} else {
	// 			fmt.Println("to delete:", file.Name())
	// 		}
	// 	} else {
	// 		fmt.Println("protected:", file.Name())
	// 	}
	// }
}

func saveFileExists(saveFile string, fileList []fs.FileInfo) bool {
	for _, file := range fileList {
		if file.Name() == saveFile {
			return true
		}
	}
	return false
}
