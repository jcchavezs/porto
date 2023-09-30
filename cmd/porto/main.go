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
	flagWriteOutputToFile := flag.Bool("w", false, "Write result to (source) file instead of stdout")
	flagListDiff := flag.Bool("l", false, "List files whose vanity import differs from porto's")
	flagSkipFiles := flag.String("skip-files", "", "Regexps of files to skip")
	flagSkipDirs := flag.String("skip-dirs", "", "Regexps of directories to skip")
	flagSkipDefaultDirs := flag.Bool("skip-dirs-use-default", true, "Use default skip directory list")
	flagIncludeInternal := flag.Bool("include-internal", false, "Include internal folders")
	restrictToFiles := flag.String("restrict-to-files", "", "Regexps of files to restrict the inspection on. It takes precedence over -skip-files and -skip-dirs")
	flag.Parse()

	baseDir := flag.Arg(0)

	if len(flag.Args()) == 0 {
		flag.Usage()
		fmt.Println(`
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

	restrictToFilesRegex, err := porto.GetRegexpList(*restrictToFiles)
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

	if len(restrictToFilesRegex) > 0 {
		opts.RestrictToFilesRegexes = restrictToFilesRegex
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
