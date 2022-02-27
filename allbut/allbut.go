package allbut

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type allbut struct {
	toDelete      []string
	toProtect     []string
	deleteEnabled bool
}

type app struct {
	p parseArgser
	s fileSanitizer
	v fileValidator
	c pathCollector
}

type parseArgser interface {
	parseArgs([]string) ([]string, bool)
}

type fileSanitizer interface {
	sanitizeFilenames([]string) []string
}

type fileValidator interface {
	validate([]string) error
}

type pathCollector interface {
	collect() ([]os.FileInfo, error)
}

func New() *app {
	p := &argParser{}
	s := &sanitizer{}
	v := &validator{}
	c := &cwdCollector{}

	return &app{
		p: p,
		s: s,
		v: v,
		c: c,
	}
}

func (a *app) Setup(args []string) (*allbut, error) {
	protectionCandidates, deleteEnabled := a.p.parseArgs(args)
	if len(protectionCandidates) == 0 {
		return &allbut{}, fmt.Errorf("%d files to protect. Provide an argument", len(protectionCandidates))
	}

	protectedFiles := a.s.sanitizeFilenames(protectionCandidates)

	err := a.v.validate(protectedFiles)
	if err != nil {
		return &allbut{}, fmt.Errorf("error validating protected files: [%s]. err %v", protectedFiles, err)
	}

	cwdFiles, err := a.c.collect()
	if err != nil {
		return &allbut{}, fmt.Errorf("error reading directory. err %v", err)
	}

	deletionCandidates, err := identifyDeletionCandidates(protectedFiles, cwdFiles)
	if err != nil {
		return &allbut{}, err
	}

	deletionTargets := a.s.sanitizeFilenames(deletionCandidates)
	return &allbut{
		toDelete:      deletionTargets,
		toProtect:     protectedFiles,
		deleteEnabled: deleteEnabled,
	}, nil
}

func (a *allbut) Run() error {
	if a.deleteEnabled {
		err := handleDeletions(a.toDelete, a.deleteEnabled)
		if err != nil {
			return fmt.Errorf("error during deletion. err: %+v", err)
		}
		return nil
	}

	// print status
	func() {
		fmt.Printf("\nProtected Files:\n")
		for _, s := range a.toProtect {
			fmt.Printf("\t%s\n", s)
		}

		fmt.Printf("\nFiles to Delete:\n")
		for _, d := range a.toDelete {
			fmt.Printf("\t%s\n", d)
		}
	}()

	// Get confirmation
	deleteConfirmation, err := getDeleteConfirmation(len(a.toDelete))
	if err != nil {
		log.Print(err)
		return err
	}

	// do the deed
	if deleteConfirmation {
		err = handleDeletions(a.toDelete, a.deleteEnabled)
		if err != nil {
			return fmt.Errorf("error during deletion. err: %+v", err)
		}
	}
	return nil
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
		if deletionEnabled {
			err := os.Remove(c)
			if err != nil {
				panic(err)
			}
			continue
		}
		fmt.Println("(use -f to really delete) MOCK deleting ", c)

	}
	return nil
}

func identifyDeletionCandidates(protectedFiles []string, filesInCwd []os.FileInfo) ([]string, error) {

	deletionCandidates := []string{}
	// Iterate and decide
	for _, fileInCwd := range filesInCwd {
		fileProtected := false
		for _, protectedFile := range protectedFiles {
			if fileInCwd.Name() == strings.ReplaceAll(protectedFile, "./", "") {
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

func PrintUsageAndExit() {
	fmt.Println("allbut deletes everything except the files you name")
	fmt.Println("usage: allbut [-f] filename1 [filename2 filename3 ...]")
	os.Exit(1)
}
