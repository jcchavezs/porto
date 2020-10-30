package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

const pathSeparator = string(os.PathSeparator)

// isGoFile checks if a file name is for a go file.
func isGoFile(filename string) bool {
	return len(filename) > 3 && strings.HasSuffix(filename, ".go")
}

// isGoTestFile checks if a file name is for a go test file.
func isGoTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

// isUnexportedDir checks if a dirname is a known unexported directory. Notice, we don't
// ignore "internal" because even when it is unexported, internally you still need to know
// the vanity import
func isUnexportedDir(dirname string) bool {
	return dirname == "testdata"
}

// addImportPath adds the vanity import path to a given go file.
func addImportPath(absFilepath string, module string) error {
	fset := token.NewFileSet()
	pf, err := parser.ParseFile(fset, absFilepath, nil, 0)
	if err != nil {
		return fmt.Errorf("failed to parse the file %q: %v", absFilepath, err)
	}
	packageName := pf.Name.String()
	if packageName == "main" { // you can't import a main package
		return nil
	}

	content, err := ioutil.ReadFile(absFilepath)
	if err != nil {
		return fmt.Errorf("failed to parse the file %q: %v", absFilepath, err)
	}

	const packageLen = 7 // = len("package")
	lenContent := len(content)
	startPackageLinePos := 0 // the position where the `package` word starts
	newLineChar := byte(10)

	if string(content[startPackageLinePos:startPackageLinePos+packageLen]) != "package" { // if first word isn't package
		startPackageLinePos = 3 // at least it started with // and a newline
		for {
			if content[startPackageLinePos-1] == newLineChar &&
				string(content[startPackageLinePos:startPackageLinePos+packageLen]) == "package" {
				// if the startPackageLinePos is a newLine and the next bytes are for package we found the position
				break
			}
			if startPackageLinePos == lenContent-packageLen {
				return fmt.Errorf("failed to find package keyword in %q", absFilepath)
			}
			startPackageLinePos++
		}
	}

	// first 1 = len(" ") as in "package " and the other 1 is for newline
	endPackageLinePos := startPackageLinePos + packageLen + 1 + 1
	for {
		if content[endPackageLinePos] == newLineChar {
			break
		}
		endPackageLinePos++
	}

	importComment := []byte(" // import \"" + module + "\"")

	newContent := []byte{}
	if startPackageLinePos != 0 {
		newContent = append(newContent, content[0:startPackageLinePos]...)
	}
	newContent = append(newContent, []byte("package "+packageName)...)
	newContent = append(newContent, importComment...)
	newContent = append(newContent, content[endPackageLinePos:]...)

	if *flagW {
		err = writeContentToFile(absFilepath, newContent)
		if err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}
	} else {
		fmt.Printf("👉 %s\n\n", absFilepath)
		fmt.Println(string(newContent))
	}

	return nil
}

// writeContentToFile writes the content in bytes to a given file.
func writeContentToFile(absFilepath string, content []byte) error {
	f, err := os.OpenFile(absFilepath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	return err
}

// findGoMod finds a go.mod file in a given directory
func findGoMod(dir string) (string, bool) {
	content, err := ioutil.ReadFile(dir + pathSeparator + "go.mod")
	if err != nil {
		return "", false
	}

	return modfile.ModulePath(content), true
}

func findAndAddVanityImportForModuleDir(absDir string, moduleName string) {
	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		log.Fatalf("failed to read the content of %q: %v", absDir, err)
	}

	for _, f := range files {
		if isDir, dirName := f.IsDir(), f.Name(); isDir {
			if isUnexportedDir(dirName) {
				continue
			} else if newModuleName, ok := findGoMod(absDir + pathSeparator + dirName); ok {
				// if folder contains go.mod we use it from now on to build the vanity import
				findAndAddVanityImportForModuleDir(absDir+pathSeparator+dirName, newModuleName)
			} else {
				// if not, we add the folder name to the vanity import
				findAndAddVanityImportForModuleDir(absDir+pathSeparator+dirName, moduleName+"/"+dirName)
			}
		} else if fileName := f.Name(); isGoFile(fileName) && !isGoTestFile(fileName) {
			if err := addImportPath(absDir+pathSeparator+fileName, moduleName); err != nil {
				log.Fatalf("failed to add vanity import path to %q: %v\n", absDir+pathSeparator+fileName, err)
			}
		}
	}
}

func findAndAddVanityImportForNonModuleDir(absDir string) {
	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		log.Fatalf("failed to read %q: %v", absDir, err)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		dirName := f.Name()
		if isUnexportedDir(dirName) {
			continue
		}

		absDirName := absDir + pathSeparator + dirName
		if moduleName, ok := findGoMod(absDirName); ok {
			findAndAddVanityImportForModuleDir(dirName, moduleName)
		} else {
			findAndAddVanityImportForNonModuleDir(absDirName)
		}
	}
}

func findAndAddVanityImportForDir(absDir string) {
	if moduleName, ok := findGoMod(absDir); ok {
		findAndAddVanityImportForModuleDir(absDir, moduleName)
		return
	}

	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		log.Fatalf("failed to read the content of %q: %v", absDir, err)
	}

	for _, f := range files {
		if !f.IsDir() {
			// we already knew this is not a Go modules folder hence we are not looking
			// for files but for directories
			continue
		}

		dirName := f.Name()
		if isUnexportedDir(dirName) {
			continue
		}

		absDirName := absDir + pathSeparator + dirName
		if moduleName, ok := findGoMod(absDirName); ok {
			findAndAddVanityImportForModuleDir(dirName, moduleName)
		} else {
			findAndAddVanityImportForNonModuleDir(absDirName)
		}
	}
}

var flagW = flag.Bool("w", false, "write result to (source) file instead of stdout")

func main() {
	flag.Parse()

	baseDir := flag.Arg(0)
	if len(flag.Args()) == 0 {
		fmt.Println(`
usage: porto [options] path

Options:
-w    write result to (source) file instead of stdout

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

	findAndAddVanityImportForDir(baseAbsDir)
}
