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
	flagW := flag.Bool("w", false, "write result to (source) file instead of stdout")
	flag.Parse()

	baseDir := flag.Arg(0)
	if len(flag.Args()) == 0 {
		fmt.Println(`
usage: porto [options] path

Options:
-w            write result to (source) file instead of stdout (default: false)

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

	err = porto.FindAndAddVanityImportForDir(baseAbsDir, porto.Options{
		WriteResultToFile: *flagW,
		GeneratedPrefixes: []string{"Code generated"},
	})
	if err != nil {
		log.Fatal(err)
	}
}
