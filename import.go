package porto

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

var (
	errGeneratedCode = errors.New("failed to add import to a generated code")
	errMainPackage   = errors.New("failed to add import to a main package")
)

// addImportPath adds the vanity import path to a given go file.
func addImportPath(absFilepath string, module string, genPrefixes []string) ([]byte, error) {
	fset := token.NewFileSet()
	pf, err := parser.ParseFile(fset, absFilepath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the file %q: %v", absFilepath, err)
	}
	packageName := pf.Name.String()
	if packageName == "main" { // you can't import a main package
		return nil, errMainPackage
	}

	content, err := ioutil.ReadFile(absFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the file %q: %v", absFilepath, err)
	}

	// 9 = len("package ") + 1 because that is the first character of the package
	startPackageLinePos := int(pf.Name.NamePos) - 9

	headerComments := string(content[0:startPackageLinePos])
	for _, genPrefix := range genPrefixes {
		if strings.Contains(headerComments, "// "+genPrefix) {
			return nil, errGeneratedCode
		}
	}

	// first 1 = len(" ") as in "package " and the other 1 is for newline
	endPackageLinePos := pf.Name.NamePos
	newLineChar := byte(10)
	for {
		// we look for new lines in case we already had comments next to the package or
		// another vanity import
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

	return newContent, nil
}

func findAndAddVanityImportForModuleDir(absDir string, moduleName string, opts Options) error {
	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		return fmt.Errorf("failed to read the content of %q: %v", absDir, err)
	}

	for _, f := range files {
		if isDir, dirName := f.IsDir(), f.Name(); isDir {
			if isUnexportedDir(dirName) {
				continue
			} else if newModuleName, ok := findGoModule(absDir + pathSeparator + dirName); ok {
				// if folder contains go.mod we use it from now on to build the vanity import
				if err := findAndAddVanityImportForModuleDir(absDir+pathSeparator+dirName, newModuleName, opts); err != nil {
					return err
				}
			} else {
				// if not, we add the folder name to the vanity import
				if err := findAndAddVanityImportForModuleDir(absDir+pathSeparator+dirName, moduleName+"/"+dirName, opts); err != nil {
					return err
				}
			}
		} else if fileName := f.Name(); isGoFile(fileName) && !isGoTestFile(fileName) {
			absFilepath := absDir + pathSeparator + fileName

			newContent, err := addImportPath(absDir+pathSeparator+fileName, moduleName, opts.GeneratedPrefixes)
			switch err {
			case nil:
				if opts.WriteResultToFile {
					err = writeContentToFile(absFilepath, newContent)
					if err != nil {
						return fmt.Errorf("failed to write file: %v", err)
					}
				} else {
					fmt.Printf("👉 %s\n\n", absFilepath)
					fmt.Println(string(newContent))
				}
			case errGeneratedCode, errMainPackage:
				continue
			default:
				return fmt.Errorf("failed to add vanity import path to %q: %v\n", absDir+pathSeparator+fileName, err)
			}
		}
	}

	return nil
}

func findAndAddVanityImportForNonModuleDir(absDir string, opts Options) error {
	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		return fmt.Errorf("failed to read %q: %v", absDir, err)
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
		if moduleName, ok := findGoModule(absDirName); ok {
			if err := findAndAddVanityImportForModuleDir(dirName, moduleName, opts); err != nil {
				return err
			}
		} else {
			if err := findAndAddVanityImportForNonModuleDir(absDirName, opts); err != nil {
				return err
			}
		}
	}
	return nil
}

type Options struct {
	// writes result to file directly
	WriteResultToFile bool
	GeneratedPrefixes []string
}

func FindAndAddVanityImportForDir(absDir string, opts Options) error {
	if moduleName, ok := findGoModule(absDir); ok {
		return findAndAddVanityImportForModuleDir(absDir, moduleName, opts)
	}

	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		return fmt.Errorf("failed to read the content of %q: %v", absDir, err)
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
		if moduleName, ok := findGoModule(absDirName); ok {
			if err := findAndAddVanityImportForModuleDir(dirName, moduleName, opts); err != nil {
				return err
			}
		} else {
			if err := findAndAddVanityImportForNonModuleDir(absDirName, opts); err != nil {
				return err
			}
		}
	}

	return nil
}
