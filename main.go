package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type DirEntry struct {
	Name  string
	IsDir bool
	Size  int64
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	err := dirPrinter(path, "", printFiles, out)
	if err != nil {
		return err
	}

	return nil
}

func dirPrinter(wd, prefix string, printFiles bool, out io.Writer) error {
	var dirContent []DirEntry //Содержимое текущей директории
	dirs, _ := ioutil.ReadDir(wd)
	//dirs, _ := os.ReadDir(wd)

	//Заполнение dirContent
	for _, dir := range dirs {

		if printFiles {
			var de DirEntry
			if dir.IsDir() {
				de = DirEntry{dir.Name(), true, -1}
			} else {
				de = DirEntry{dir.Name(), false, dir.Size()}
			}

			dirContent = append(dirContent, de)

		} else {
			if dir.IsDir() {
				var de DirEntry
				de = DirEntry{dir.Name(), true, -1}
				dirContent = append(dirContent, de)
			}
		}

	}

	//Сортировка содержимого директории по имени
	sort.Slice(dirContent, func(i, j int) bool {
		return dirContent[i].Name < dirContent[j].Name
	})

	//Вывод (тестовый)
	for i, dir := range dirContent {

		if dir.IsDir {
			var newPrefix string
			newWd := wd + string(os.PathSeparator) + dir.Name

			if i == len(dirContent)-1 {
				newPrefix = prefix + "\t"
				fmt.Fprintln(out, prefix+"└───"+dir.Name)
			} else {
				newPrefix = prefix + "│\t"
				fmt.Fprintln(out, prefix+"├───"+dir.Name)
			}
			err := dirPrinter(newWd, newPrefix, printFiles, out)
			if err != nil {
				return err
			}
		} else {
			if printFiles {
				if dir.Size == 0 {
					if i == len(dirContent)-1 {
						fmt.Fprintln(out, prefix+"└───"+dir.Name+" ("+"empty"+")")
					} else {
						fmt.Fprintln(out, prefix+"├───"+dir.Name+" ("+"empty"+")")
					}
				} else {
					if i == len(dirContent)-1 {
						fmt.Fprintln(out, prefix+"└───"+dir.Name+" ("+strconv.FormatInt(dir.Size, 10)+"b)")
					} else {
						fmt.Fprintln(out, prefix+"├───"+dir.Name+" ("+strconv.FormatInt(dir.Size, 10)+"b)")
					}
				}
			}

		}
	}

	return nil
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
