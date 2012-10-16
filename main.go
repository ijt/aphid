package main;

import (
	"bufio"
	"fmt"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

import "github.com/kylelemons/go-gypsy/yaml"

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
	defaultUrl := "https://raw.github.com/ijt/aphid/config/config.yaml"
	defaultPrefix := "  " + Bold + FgCyan + "[aphid]" + Reset
	configUrl := flag.String("c", defaultUrl,
				 "URL of config file in YAML format")
	prefix := flag.String("p", defaultPrefix, "Prefix for aphid messages")
	flag.Parse()

	config, err := parseConfig(fetch(*configUrl), *configUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	compileRegexes(config)
	addHelp(bufio.NewReader(os.Stdin), config, *prefix)
}

// fetch gets the contents at a given URL. The URL can point to a local file.
// Errors terminate.
func fetch(url string) []byte {
	// Make a client that can load files if given a file:// URL.
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c := &http.Client{Transport: t}

	// Download the config file from a well-known location.
	resp, err := c.Get(url)
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

	return body
}

func parseConfig(body []byte, url string) (conf *Config, err error) {
	defer func() {
		if e := recover(); e != nil {
			conf = nil      // Clear return value.
			body := string(body)
			if len(body) > 200 {
				body = body[:200] + "..."
			}
			err = fmt.Errorf("Failed to parse config file at %s:\n\n%s",
					 url, body)
		}
	}()

	// Parse the YAML file.
	yamlFile := yaml.Config(string(body))
	count, e := yamlFile.Count("line_rules")
	if e != nil {
		err = fmt.Errorf("Could not get line_rules config section.")
		return
	}

	// Extract the line rules.
	conf = &Config {}
	conf.Line_rules = make([]*LineRule, count)
	for i, _ := range conf.Line_rules {
		pattern, e := yamlFile.Get(fmt.Sprintf("line_rules[%d].pattern", i))
		if e != nil {
			return nil, e
		}

		message, e := yamlFile.Get(fmt.Sprintf("line_rules[%d].message", i))
		if e != nil {
			return nil, e
		}

		conf.Line_rules[i] = &LineRule{ pattern, nil, message }
	}

	return conf, nil
}

func compileRegexes(conf *Config) *Config {
	for _, rule := range conf.Line_rules {
		rule.patternRx = regexp.MustCompile(".*" + rule.Pattern + ".*")
		if rule.patternRx == nil {
			fmt.Println("Pattern failed to compile:", rule.Pattern)
			os.Exit(1)
		}
	}
	return conf
}

// addHelp reads lines from a reader and prints them to stdout. It interjects
// helpful messages when those lines match patterns given by the config.
func addHelp(reader *bufio.Reader, conf *Config, prefix string) {
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
				msg := rule.patternRx.ReplaceAllString(
					     line, rule.Message)
				fmt.Print(prefix, " ", msg)
			}
		}
	}
}

// https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/99MKtEkvQ2c
const (
        Reset = "\x1b[0m"
        Bold = "\x1b[1m"
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

