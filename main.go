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
	Message_prefix string
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

	if conf.Message_prefix == "" {
		conf.Message_prefix = "  " + FgCyan + "[catkin_sleuth]" + Reset
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
				fmt.Print(conf.Message_prefix, " ", rule.Message)	
			}
		}
	}
}

const (
        Reset = "\x1b[0m"
        Bright = "\x1b[1m"
        Dim = "\x1b[2m"
        Underscore = "\x1b[4m"
        Blink = "\x1b[5m"
        Reverse = "\x1b[7m"
        Hidden = "\x1b[8m"

        FgBlack = "\x1b[30m"
        FgRed = "\x1b[31m"
        FgGreen = "\x1b[32m"
        FgYellow = "\x1b[33m"
        FgBlue = "\x1b[34m"
        FgMagenta = "\x1b[35m"
        FgCyan = "\x1b[36m"
        FgWhite = "\x1b[37m"

        BgBlack = "\x1b[40m"
        BgRed = "\x1b[41m"
        BgGreen = "\x1b[42m"
        BgYellow = "\x1b[43m"
        BgBlue = "\x1b[44m"
        BgMagenta = "\x1b[45m"
        BgCyan = "\x1b[46m"
        BgWhite = "\x1b[47m"
)

