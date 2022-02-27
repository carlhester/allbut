package allbut

type sanitizer struct{}

func (s *sanitizer) sanitizeFilenames(files []string) []string {
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
