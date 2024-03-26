package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dfsTree(out, path, printFiles, 0, 0)
}

func dfsTree(out io.Writer, path string, printFiles bool, lvl int, prevNotLastCnt int) error {
	openedEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var entries []os.DirEntry
	if printFiles {
		entries = make([]os.DirEntry, len(openedEntries))
		copy(entries, openedEntries)
	} else {
		for i := range openedEntries {
			info, err := openedEntries[i].Info()
			if err != nil {
				return err
			}

			if info.IsDir() {
				entries = append(entries, openedEntries[i])
			}
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for i, e := range entries {
		isLast := len(entries)-1 == i
		err := printGraphics(out, lvl, isLast, printFiles, prevNotLastCnt, e)
		if err != nil {
			return err
		}

		if e.IsDir() {
			newPath := path + string(os.PathSeparator) + e.Name()
			currNotLastCnt := prevNotLastCnt
			if !isLast {
				currNotLastCnt++
			}
			err = dfsTree(out, newPath, printFiles, lvl+1, currNotLastCnt)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func printGraphics(out io.Writer, lvl int, isLast bool, printFiles bool, prevNotLastCnt int, e os.DirEntry) error {
	if !e.IsDir() && !printFiles {
		return nil
	}

	for i := 0; i < lvl; i++ {
		if i < prevNotLastCnt {
			_, err := out.Write([]byte("│"))
			if err != nil {
				return err
			}
		}

		_, err := out.Write([]byte("\t"))
		if err != nil {
			return err
		}
	}

	var info string
	var err error
	if e.IsDir() {
		info = e.Name()
	} else {
		info, err = getFileInfo(e)
		if err != nil {
			return err
		}
	}

	if isLast {
		str := "└───" + info + "\n"
		_, err = out.Write([]byte(str))
	} else {
		str := "├───" + info + "\n"
		_, err = out.Write([]byte(str))
	}
	return err
}

func getFileInfo(e os.DirEntry) (string, error) {
	info, err := e.Info()
	if err != nil {
		return "", err
	}

	var size string
	if info.Size() > 0 {
		size = fmt.Sprintf("%db", info.Size())
	} else {
		size = "empty"
	}
	return fmt.Sprintf("%s (%s)", e.Name(), size), nil
}
