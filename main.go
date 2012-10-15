package main;

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	sections := make(chan []string)
	go sectionize(bufio.NewReader(os.Stdin), sections)
	sections2 := make(chan []string)
	go addHelpToSections(sections, sections2)
	printSections(sections2)
}

func sectionize(reader *bufio.Reader, sections chan []string) {
	var section []string;
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRightFunc(line, unicode.IsSpace)
		if len(line) > 1 && line[0] != ' ' {
			if section != nil {
				sections <- section
			}
			a := []string {line}
			section = a[:]
		} else {
			section = append(section, line)
		}
	}
	close(sections)
}

func addHelpToSections(sections chan []string, sections2 chan []string) {
	for s := range(sections) {
		pattern := "Could not find a configuration file for package"
		if len(s) >= 2 && strings.Contains(s[1], pattern) {
			s = append(s, "")
			s = append(s, "  [cmake_sleuth] Did you install this package?")
			s = append(s, "  [cmake_sleuth] Did you run build/buildspace/setup.sh?")
			s = append(s, "")
		}

		sections2 <- s
	}
	close(sections2)
}

func printSections(sections chan []string) {
	for s := range(sections) {
		for _, line := range(s) {
			fmt.Println(line)
		}
	}
}

