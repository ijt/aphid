package main;

import (
	"bufio"
	"fmt"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

import "github.com/kylelemons/go-gypsy/yaml"

type Config struct {
	lineRules []*LineRule
}

type LineRule struct {
	pattern string
	patternRx *regexp.Regexp
	message string
}

func main() {
	defaultUrl := "https://raw.github.com/ijt/aphid/config/config.yaml"
	defaultPrefix := "  " + Bold + FgCyan + "[aphid]" + Reset
	configUrl := flag.String("c", defaultUrl,
				 "URL of config file in YAML format")
	prefix := flag.String("p", defaultPrefix, "Prefix for aphid messages")
	strict := flag.Bool("s", false, "Be strict when parsing patterns")
	flag.Parse()

	config, err := parseConfig(fetch(*configUrl), *configUrl)
	if err != nil {
		log.Fatalln(err)
	}
	err = compileRegexes(config)
	if err != nil && *strict {
		os.Exit(1)
	}
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
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func parseConfig(body []byte, url string) (conf *Config, err error) {
	defer func() {
		if e := recover(); e != nil {
			// Clear return value.
			conf = nil
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
	count, e := yamlFile.Count("lineRules")
	if e != nil {
		err = fmt.Errorf("Could not get lineRules config section.")
		return
	}

	// Extract the line rules.
	conf = &Config {}
	conf.lineRules = make([]*LineRule, count)
	for i, _ := range conf.lineRules {
		pattern, e := yamlFile.Get(fmt.Sprintf("lineRules[%d].pattern", i))
		if e != nil {
			return nil, e
		}

		message, e := yamlFile.Get(fmt.Sprintf("lineRules[%d].message", i))
		if e != nil {
			return nil, e
		}

		conf.lineRules[i] = &LineRule{ pattern, nil, message }
	}

	return conf, nil
}

// compileRegexes tries to compile all the regexes in the config.
// If any of them fail, it returns the first error it found.
func compileRegexes(conf *Config) error {
	retErr := error(nil)
	for _, rule := range conf.lineRules {
		rx, err := regexp.Compile(".*" + rule.pattern + ".*")
		if err != nil {
			log.Println(err)
			if retErr != nil {
				retErr = err
			}
		} else {
			rule.patternRx = rx
		}
	}
	return retErr
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
		for _, rule := range(conf.lineRules) {
			matched := rule.patternRx.MatchString(line)
			if matched {
				msg := rule.patternRx.ReplaceAllString(
					     line, rule.message)
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

