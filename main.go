package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// Parse input
	// Strip the first argument which is the executable itself
	args := os.Args[1:]

	protectionCandidates, deleteFlagEnabled := parseArgs(args)
	if len(protectionCandidates) == 0 {
		printUsageAndExit()
	}

	// Sanitize files
	sanitizedFiles := sanitizeFileNames(protectionCandidates)

	// validate input files
	protectedFiles, err := validateFiles(sanitizedFiles)
	if err != nil {
		log.Printf("error validating protected files: [%s]. err %v", protectedFiles, err)
		printUsageAndExit()
	}

	// collect files from current directory
	cwdFiles, err := ioutil.ReadDir("./")
	if err != nil {
		log.Printf("error reading directory. err %v", err)
		printUsageAndExit()
	}

	// process files
	deletionCandidates, err := identifyDeletionCandidates(protectedFiles, cwdFiles)
	if err != nil {
		log.Print(err)
		printUsageAndExit()
	}

	deletionTargets := sanitizeFileNames(deletionCandidates)

	// print status
	func() {
		fmt.Println("Protected Files:")
		for _, s := range protectedFiles {
			fmt.Printf("\t%s\n", s)
		}

		fmt.Printf("\nFiles to Delete:\n")
		for _, d := range deletionTargets {
			fmt.Printf("\t%s\n", d)
		}
	}()

	// Get confirmation
	deleteConfirmation, err := getDeleteConfirmation(len(deletionCandidates))
	if err != nil {
		log.Print(err)
		printUsageAndExit()
	}

	// do the deed
	if deleteConfirmation {
		err = handleDeletions(deletionCandidates, deleteFlagEnabled)
		if err != nil {
			log.Print(err)
			printUsageAndExit()
		}
	}
}

func removeInvalidChars(i string) string {
	badChars := "[],~\\!@#$%^&*(){}'<>?;:=+'"
	toStrip := strings.Split(badChars, "")
	for _, s := range toStrip {
		i = strings.Replace(i, string(s), "", -1)
	}
	return i
}

func addDotSlashPrefix(i string) string {
	s := strings.Replace(i, "./", "", -1)
	return fmt.Sprintf("./%s", s)
}

func sanitizeFileNames(files []string) []string {
	r := []string{}
	stripped := []string{}

	for _, f := range files {
		p := addDotSlashPrefix(f)
		stripped = append(stripped, p)
	}

	for _, f := range stripped {
		r = append(r, removeInvalidChars(f))
	}

	return r
}

func getDeleteConfirmation(count int) (bool, error) {
	fmt.Printf("\n%d files will be deleted. Type [y] to proceed: ", count)
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return false, fmt.Errorf("error getting confirmation. err: %s", err)
	}
	if string(char) == "y" {
		return true, nil
	}
	return false, nil
}

func handleDeletions(candidates []string, deletionEnabled bool) error {
	for _, c := range candidates {
		if !deletionEnabled {
			fmt.Println("MOCK deleting ... ", c)
		}
	}
	return nil
}

func identifyDeletionCandidates(protectedFiles []string, filesInCwd []os.FileInfo) ([]string, error) {
	deletionCandidates := []string{}
	// Iterate and decide
	for _, fileInCwd := range filesInCwd {
		fileProtected := false
		for _, protectedFile := range protectedFiles {
			if fileInCwd.Name() == protectedFile {
				fileProtected = true
				continue
			}
		}
		if !fileProtected {
			deletionCandidates = append(deletionCandidates, fileInCwd.Name())
		}
	}
	return deletionCandidates, nil
}

func validateFiles(files []string) ([]string, error) {
	for _, file := range files {
		f, err := os.Stat(file)

		// Check that file exists
		if err != nil {
			return []string{}, fmt.Errorf("unable to read file: %s", f)
		}
		// Check that file is not a directory
		if f.IsDir() {
			return []string{}, fmt.Errorf("%s is a directory, not a plain file", f.Name())
		}
	}
	return files, nil
}

func printUsageAndExit() {
	fmt.Println("allbut deletes everything except the files you name")
	fmt.Println("usage: allbut [-f] filename1 [filename2 filename3 ...]")
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
