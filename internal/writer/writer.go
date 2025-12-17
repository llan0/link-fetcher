package writer

import (
	"os"
	"strings"
)

func WriteLinks(filename string, links []string) error {
	if len(links) == 0 {
		return nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var batch []string

	for _, link := range links {
		batch = append(batch, link)

		if len(batch) == 10 {
			line := strings.Join(batch, ",")
			if _, err := f.WriteString(line + "\n\n"); err != nil {
				return err
			}
			batch = nil // Reset
		}
	}

	if len(batch) > 0 {
		line := strings.Join(batch, ",")
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}
