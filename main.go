package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	help = flag.Bool("help", false, "Show help")
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:
    process-scanner pattern [patterns...]
    process-scanner --help
`)
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func buildRegexpPattern(patterns []string) string {
	quotedPatterns := make([]string, 0)
	for _, pattern := range patterns {
		quotedPatterns = append(quotedPatterns, regexp.QuoteMeta(pattern))
	}

	return strings.Join(quotedPatterns, "|")
}

func scanForProcs(patterns []string) error {
	patternRegexp := regexp.MustCompile(buildRegexpPattern(patterns))

	// Compile the regular expression used to match process directories.
	procRegexp := regexp.MustCompile("^\\d+$")

	// Open the root directory of the proc filesystem.
	procfs, err := os.Open("/proc")
	if err != nil {
		return fmt.Errorf("unable to open /proc: %s", err.Error())
	}

	// List files in the proc filesystem.
	dirContents, err := procfs.Readdir(0)
	if err != nil {
		return fmt.Errorf("unable to list files in /proc: %s", err.Error())
	}

	// Just print the directory contents for now.
	for _, fileInfo := range dirContents {

		// We only care about files that are dirctories and whose names look like process IDs.
		if !fileInfo.IsDir() || !procRegexp.MatchString(fileInfo.Name()) {
			continue
		}

		// Open the command line file; skip this process if it has no command line.
		cmdlinePath := path.Join(procfs.Name(), fileInfo.Name(), "cmdline")
		cmdlineFile, err := os.Open(cmdlinePath)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("unable to open %s: %s", cmdlinePath, err.Error())
		}

		// Extract the arguments from the command-line file.
		cmdline, err := ioutil.ReadAll(cmdlineFile)
		if err != nil {
			return fmt.Errorf("unable to read %s: %s", cmdlinePath, err.Error())
		}
		args := strings.Split(string(cmdline), "\000")

		// Print the command line if the first command itself matches one of the patterns.
		if patternRegexp.Match([]byte(args[0])) {
			fmt.Println(strings.Join(args, " "))
		}
	}
	return nil
}

func main() {
	if *help {
		usage()
		os.Exit(0)
	}

	// We have to have at least one pattern to search for.
	if len(os.Args) <= 1 {
		usage()
		os.Exit(1)
	}

	// Scan for processes matching one of the patterns.
	_ = scanForProcs(os.Args[1:])
}
