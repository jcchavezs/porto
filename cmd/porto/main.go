package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jcchavezs/porto"
)

func getRegexpList(flagVal string) ([]*regexp.Regexp, error) {
	var regexes []*regexp.Regexp
	if len(flagVal) > 0 {
		for _, sfrp := range strings.Split(flagVal, ",") {
			sfr, err := regexp.Compile(sfrp)
			if err != nil {
				return nil, fmt.Errorf("failed to compile regex %q: %w", sfr, err)
			}
			regexes = append(regexes, sfr)
		}
	}

	return regexes, nil
}

func main() {
	flagWriteOutputToFile := flag.Bool("w", false, "write result to (source) file instead of stdout")
	flagListDiff := flag.Bool("l", false, "list files whose vanity import differs from porto's")
	flagSkipFiles := flag.String("skip-files", "", "Regexps of files to skip")
	flagIncludeInternal := flag.Bool("include-internal", false, "include internal folders")
	flag.Parse()

	baseDir := flag.Arg(0)
	if len(flag.Args()) == 0 {
		fmt.Println(`
usage: porto [options] <target-path>

Options:
-w                  Write result to (source) file instead of stdout (default: false)
-l                  List files whose vanity import differs from porto's (default: false)
--skip-files        Regexps of files to skip
--include-internal  Include internal folders

Examples:

Add import path to a folder
    $ porto -w ./myproject
		`)
		os.Exit(0)
	}

	baseAbsDir, err := filepath.Abs(baseDir)
	if err != nil {
		log.Fatalf("failed to resolve base absolute path for target path %q: %v", baseDir, err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to resolve base absolute path for current working dir: %v", err)
	}

	skipFilesRegex, err := getRegexpList(*flagSkipFiles)
	if err != nil {
		log.Fatalf("failed to build files regexes: %v", err)
	}

	diffCount, err := porto.FindAndAddVanityImportForDir(workingDir, baseAbsDir, porto.Options{
		WriteResultToFile: *flagWriteOutputToFile,
		ListDiffFiles:     *flagListDiff,
		SkipFilesRegexes:  skipFilesRegex,
		IncludeInternal:   *flagIncludeInternal,
	})
	if err != nil {
		log.Fatal(err)
	}

	if *flagListDiff && diffCount > 0 {
		os.Exit(2)
	}
}
