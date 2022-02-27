package allbut

type argParser struct{}

func (p *argParser) parseArgs(args []string) ([]string, bool) {
	results := []string{}
	deletionEnabled := false

	for _, file := range args {
		if file == "-f" {
			deletionEnabled = true
			continue
		}
		results = append(results, file)
	}
	return results, deletionEnabled
}
