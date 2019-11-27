package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	//"strings"
)

const middleIcon string = "├───"
const lastIcon string = "└───"

var filesToSkip = [...]string{".git", ".gitignore", "dockerfile", "hw1.md"}

type FolderElement struct {
	name     string
	isFolder bool
	file     os.FileInfo
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := dirTree2(out, path, printFiles, 0, "")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func dirTree2(out io.Writer, path string, printFiles bool, level int, prefix string) error {

	sortedFolderElements := getSortedDirElements(path, printFiles)
	sLen := len(sortedFolderElements) - 1

	for idx, f := range sortedFolderElements {
		isLast := isLast(sLen, idx)
		newPrefix := getPrefix(isLast, prefix)
		if f.isFolder {
			txt := prefix + getIcon(isLast) + f.name + "\n"
			bytes := []byte(txt)
			_, err := out.Write(bytes)
			if err != nil {
				log.Fatal(err)
			}

			err = dirTree2(out, filepath.Join(path, f.name), printFiles, level+1, newPrefix)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			txt := prefix + getIcon(isLast) + f.name + StringifyByteCountBinary(f.file.Size()) + "\n"
			bytes := []byte(txt)
			_, err := out.Write(bytes)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	return nil
}

func getPrefix(isLast bool, prefixPre string) (prefix string) {
	if isLast {
		prefix = prefixPre + "\t"
	} else {
		prefix = prefixPre + "│\t" +
			""
	}

	return prefix
}

func sortFolderElementsByName(elements []FolderElement) []FolderElement {
	sort.Slice(elements, func(i, j int) bool { return elements[i].name < elements[j].name })
	return elements
}

func getSortedDirElements(path string, appendFiles bool) []FolderElement {
	elements := listDir(path, appendFiles)
	return sortFolderElementsByName(elements)
	//return sorted
}

func listDir(path string, appendFiles bool) (folderElements []FolderElement) {
	folderContent, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range folderContent {

		if flag := skipElem(f); flag {
			continue
		}

		//folderElements = append(folderElements, FolderElement{"yzq", true, f})

		if dirFlag := f.IsDir(); dirFlag {
			folderElements = append(folderElements, FolderElement{f.Name(), true, f})
		} else if appendFiles {
			folderElements = append(folderElements, FolderElement{f.Name(), false, f})
		}

	}

	return folderElements
}

func skipElem(f os.FileInfo) bool {
	for _, a := range filesToSkip {
		if a == f.Name() {
			return true
		}
	}
	return false
}

func StringifyByteCountBinary(b int64) string {
	if b == 0 {
		return " (empty)"
	}
	return fmt.Sprintf(" (%db)", b)
}

func getIcon(isLast bool) string {
	if isLast {
		return lastIcon
	}

	return middleIcon
}

func isLast(sliceLen int, currentPos int) bool {
	if sliceLen == currentPos {
		return true
	}

	return false
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
