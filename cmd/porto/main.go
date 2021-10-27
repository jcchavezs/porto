package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jcchavezs/porto"
)

func main() {
	flagWriteOutputToFile := flag.Bool("w", false, "write result to (source) file instead of stdout")
	flagListDiff := flag.Bool("l", false, "list files whose vanity import differs from porto's")
	flagSkipFiles := flag.String("skip-files", "", "Regexps of files to skip")
	flagSkipDirs := flag.String("skip-dirs", "", "Regexps of directories to skip")
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
--skip-dirs         Regexps of directories to skip
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

	skipFilesRegex, err := porto.GetRegexpList(*flagSkipFiles)
	if err != nil {
		log.Fatalf("failed to build files regexes: %v", err)
	}

	skipDirsRegex, err := porto.GetRegexpList(*flagSkipDirs)
	if err != nil {
		log.Fatalf("failed to build directories regexes: %v", err)
	}

	diffCount, err := porto.FindAndAddVanityImportForDir(workingDir, baseAbsDir, porto.Options{
		WriteResultToFile: *flagWriteOutputToFile,
		ListDiffFiles:     *flagListDiff,
		SkipFilesRegexes:  skipFilesRegex,
		SkipDirsRegexes:   skipDirsRegex,
		IncludeInternal:   *flagIncludeInternal,
	})
	if err != nil {
		log.Fatal(err)
	}

	if *flagListDiff && diffCount > 0 {
		os.Exit(2)
	}
}
