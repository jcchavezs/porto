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
	flagWrite := flag.Bool("w", false, "write result to (source) file instead of stdout")
	flagList := flag.Bool("l", false, "list files whose vanity import differs from porto's")
	flag.Parse()

	baseDir := flag.Arg(0)
	if len(flag.Args()) == 0 {
		fmt.Println(`
usage: porto [options] path

Options:
-w            write result to (source) file instead of stdout (default: false)
-l            list files whose vanity import differs from porto's (default: false)

Examples:

Add import path to a folder
    $ porto -w myproject
		`)
		os.Exit(0)
	}

	baseAbsDir, err := filepath.Abs(baseDir)
	if err != nil {
		log.Fatalf("failed to resolve base absolute path for %q: %v", baseDir, err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to resolve base absolute path for %q: %v", baseDir, err)
	}

	err = porto.FindAndAddVanityImportForDir(workingDir, baseAbsDir, porto.Options{
		WriteResultToFile: *flagWrite,
		ListDiffFiles:     *flagList,
		GeneratedPrefixes: []string{"Code generated"},
	})
	if err != nil {
		log.Fatal(err)
	}
}
