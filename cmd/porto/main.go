package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jcchavezs/porto"
)

func main() {
	flagWriteOutputToFile := flag.Bool("w", false, "write result to (source) file instead of stdout")
	flagListDiff := flag.Bool("l", false, "list files whose vanity import differs from porto's")
	flagSkipFiles := flag.String("skip-files", "", "Regexps of files to skip")
	flagSkipDirs := flag.String("skip-dirs", "", "Regexps of directories to skip")
	flagSkipDefaultDirs := flag.Bool("skip-dirs-use-default", true, "use default skip directory list")
	flagIncludeInternal := flag.Bool("include-internal", false, "include internal folders")
	flagIncludeFiles := flag.String("include-files", "", "Regexps of files to include")
	flag.Parse()

	baseDir := flag.Arg(0)
	if len(flag.Args()) == 0 {
		fmt.Println(`
usage: porto [options] <target-path>

Options:
-w                       Write result to (source) file instead of stdout (default: false)
-l                       List files whose vanity import differs from porto's (default: false)
--skip-files             Regexps of files to skip
--skip-dirs              Regexps of directories to skip
--skip-dirs-use-default  Use default skip directory list (default: true)
--include-internal       Include internal folders
--include-files          Regexps of files to include. It takes precedence over --skip-files

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

	skipFilesRegex, err := porto.GetRegexpList(*flagSkipFiles)
	if err != nil {
		log.Fatalf("failed to build files regexes to exclude: %v", err)
	}

	includeFilesRegex, err := porto.GetRegexpList(*flagIncludeFiles)
	if err != nil {
		log.Fatalf("failed to build files regexes to include: %v", err)
	}

	var skipDirsRegex []*regexp.Regexp
	if *flagSkipDefaultDirs {
		skipDirsRegex = append(skipDirsRegex, porto.StdExcludeDirRegexps...)
	}
	userSkipDirsRegex, err := porto.GetRegexpList(*flagSkipDirs)
	if err != nil {
		log.Fatalf("failed to build directories regexes: %v", err)
	}
	skipDirsRegex = append(skipDirsRegex, userSkipDirsRegex...)

	opts := porto.Options{
		WriteResultToFile: *flagWriteOutputToFile,
		ListDiffFiles:     *flagListDiff,
		IncludeInternal:   *flagIncludeInternal,
	}

	if len(includeFilesRegex) > 0 {
		opts.IncludeFilesRegexes = includeFilesRegex
	} else {
		opts.SkipFilesRegexes = skipFilesRegex
		opts.SkipDirsRegexes = skipDirsRegex
	}

	diffCount, err := porto.FindAndAddVanityImportForDir(workingDir, baseAbsDir, opts)
	if err != nil {
		log.Fatal(err)
	}

	if *flagListDiff && diffCount > 0 {
		os.Exit(2)
	}
}
