package main;

import (
	"bufio"
	"fmt"
	"flag"
	"io/ioutil"
	"launchpad.net/goyaml"
	"net/http"
	"os"
	"regexp"
)

type Config struct {
	Line_rules []*LineRule
}

type LineRule struct {
	Pattern string
	patternRx *regexp.Regexp
	Message string
}

func main() {
	defUrl := "https://raw.github.com/ijt/catkin_sleuth/config/config.yaml"
	configUrl := flag.String("-c", defUrl,
				 "URL of config file in YAML format")
	flag.Parse()

	config := loadConfig(*configUrl)
	addHelp(bufio.NewReader(os.Stdin), config)
}

func loadConfig(url string) *Config {
	// Download the config file from a well-known location
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Parse the config file
	conf := &Config {}
	err = goyaml.Unmarshal(body, conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Compile the regexes
	for _, rule := range conf.Line_rules {
		rule.patternRx = regexp.MustCompile(rule.Pattern)
		if rule.patternRx == nil {
			fmt.Println("Pattern failed to compile:", rule.Pattern)
			os.Exit(1)
		}
	}

	return conf
}

func addHelp(reader *bufio.Reader, conf *Config) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if line == "" {
			fmt.Fprintln(os.Stderr, "Line is null. Ending.")
			os.Exit(1)
		}
		fmt.Print(line)
		for _, rule := range(conf.Line_rules) {
			matched := rule.patternRx.MatchString(line)
			if matched {
				// FIXME: Substitute positional references from
				// regexp
				fmt.Println(rule.Message)	
			}
		}
	}
}

