package gosloc

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type GoSLOC struct {
	Dir     string
	Files   []fs.DirEntry
	Results map[string]interface{}
	Total   int
}

func (gosloc *GoSLOC) SaveOrDisplay(fp string, isDisp bool) error {
	var file *os.File
	var err error
	if fp != "" {
		if path.Ext(fp) != ".txt" {
			fp += ".txt"
		}
		file, err = os.Create(fp)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	for fp, lines := range gosloc.Results {
		content := fmt.Sprintf("%d %s %s\n", lines, strings.Repeat("-", 30), fp)
		if isDisp {
			fmt.Print(content)
		}
		if fp != "" {
			file.WriteString(content)
		}
	}

	// add total
	totalCont := fmt.Sprintf("\n%d %s %s", gosloc.Total, strings.Repeat("-", 30), "Total")
	fmt.Println(totalCont)
	file.WriteString(totalCont)

	return nil
}

func contains(arr []string, v string) bool {
	if len(arr) == 0 {
		return true
	} else {
		for _, a := range arr {
			if a == v {
				return true
			}
		}
	}
	return false
}

func (gosloc *GoSLOC) Read(dir string) error {
	var err error
	gosloc.Dir = dir
	gosloc.Files, err = os.ReadDir(dir)
	if err != nil {
		return err
	}
	return nil
}

// handle recursive
func (gosloc *GoSLOC) Process(e string) error {
	var exts []string
	if e != "" {
		exts = strings.Split(e, ",")
	}

	gosloc.Results = make(map[string]interface{})
	for _, file := range gosloc.Files {
		if !file.IsDir() {
			if contains(exts, filepath.Ext(file.Name())) {
				filePath := filepath.Join(gosloc.Dir, file.Name())
				lines, err := getNumOfLines(filePath)
				if err != nil {
					return err
				}
				gosloc.Total += lines
				gosloc.Results[filePath] = lines
			}
		}
	}
	return nil
}

func getNumOfLines(filePath string) (int, error) {
	var err error
	var output bytes.Buffer

	catCmd := exec.Command("cat", filePath)
	wcCmd := exec.Command("wc", "-l")

	wcCmd.Stdout = &output
	wcCmd.Stdin, err = catCmd.StdoutPipe()
	if err != nil {
		return 0, err
	}

	if err = wcCmd.Start(); err != nil { // runs but doesn't wait to finish
		return 0, err
	}

	if err = catCmd.Run(); err != nil { // runs, waits for the output to generated
		return 0, err
	}

	if err = wcCmd.Wait(); err != nil {
		return 0, err
	}

	lines, err := strconv.Atoi(strings.Trim(output.String(), "\n"))
	if err != nil {
		return 0, err
	}
	return lines, nil
}
