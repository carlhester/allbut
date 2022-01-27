package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// Parse input
	saveFiles, deletionEnabled := parseArgs(os.Args)

	// validate input files
	for _, saveFile := range saveFiles {
		err := validateFile(saveFile)
		if err != nil {
			log.Printf("error validating file %s. err %v", saveFile, err)
			printUsageAndExit()
		}
	}

	// collect files from current directory
	filesInCwd, err := ioutil.ReadDir("./")
	if err != nil {
		log.Printf("error reading directory. err %v", err)
		printUsageAndExit()
	}

	// process files
	deletionCandidates, err := identifyDeletionCandidates(saveFiles, filesInCwd)
	if err != nil {
		log.Print(err)
		printUsageAndExit()
	}

	confirm, err := printStatusAndVerify(saveFiles, deletionCandidates)
	if err != nil {
		log.Print(err)
		printUsageAndExit()
	}

	if confirm {
		err = handleDeletions(deletionCandidates, deletionEnabled)
		if err != nil {
			log.Print(err)
			printUsageAndExit()
		}
	}
}

func printStatusAndVerify(s []string, d []fs.FileInfo) (bool, error) {
	fmt.Println("Protected Files:")
	for _, ss := range s {
		fmt.Printf("\t%s\n", ss)
	}

	fmt.Printf("\nFiles to Delete:\n")
	for _, dd := range d {
		fmt.Printf("\t%s\n", dd.Name())
	}

	fmt.Printf("\n %d files will be deleted. Type [y] to proceed: ", len(d))
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return false, fmt.Errorf("Error getting confirmation. err: %s", err)
	}
	if string(char) == "y" {
		return true, nil
	}
	return false, nil
}

func handleDeletions(candidates []os.FileInfo, deletionEnabled bool) error {
	for _, c := range candidates {
		if !deletionEnabled {
			fmt.Println("MOCK deleting ... ", c.Name())
		}
	}
	return nil
}

func identifyDeletionCandidates(saveFiles []string, filesInCwd []os.FileInfo) ([]os.FileInfo, error) {
	deletionCandidates := []os.FileInfo{}
	// Iterate and decide
	for _, fileInCwd := range filesInCwd {
		fileProtected := false
		for _, saveFile := range saveFiles {
			if fileInCwd.Name() == saveFile {
				fileProtected = true
				continue
			}
		}
		if !fileProtected {
			deletionCandidates = append(deletionCandidates, fileInCwd)
		}
	}
	return deletionCandidates, nil
}

func validateFile(file string) error {
	f, err := os.Stat(file)
	// Check that file exists
	if err != nil {
		return fmt.Errorf("unable to read file: %s", f)
	}
	// Check that file is not a directory
	if f.IsDir() {
		return fmt.Errorf("%s is a directory, not a plain file", f.Name())
	}
	return nil

}

func printUsageAndExit() {
	fmt.Println("usage: allbut [-f] filename")
	os.Exit(1)
}

func parseArgs(args []string) ([]string, bool) {
	results := []string{}
	deletionEnabled := false

	if len(args) < 2 {
		return []string{}, false
	}
	for _, file := range args {
		if file == "-f" {
			deletionEnabled = true
			continue
		}
		results = append(results, file)
	}
	return results, deletionEnabled
}
