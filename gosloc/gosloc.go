package gosloc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type GoSLOC struct {
	Dir       string
	FilePaths map[string]interface{}
	Total     int
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

	for fp, lines := range gosloc.FilePaths {
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
	if isDisp {
		fmt.Println(totalCont)
	}
	if fp != "" {
		file.WriteString(totalCont)
	}
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

func (gosloc *GoSLOC) Recursive(dir string, rec bool, exts []string) error {
	var err error
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && contains(exts, filepath.Ext(file.Name())) {
			gosloc.FilePaths[filepath.Join(dir, file.Name())] = 0
		}
		if rec && file.IsDir() {
			err = gosloc.Recursive(filepath.Join(dir, file.Name()), rec, exts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gosloc *GoSLOC) Read(dir string, rec bool, e string) error {
	var exts []string
	if e != "" {
		exts = strings.Split(e, ",")
	}
	gosloc.Dir = dir
	err := gosloc.Recursive(dir, rec, exts)
	if err != nil {
		return err
	}
	return nil
}

// handle recursive
func (gosloc *GoSLOC) Process(e string) error {
	for fp := range gosloc.FilePaths {
		lines, err := getNumOfLines(fp)
		if err != nil {
			return err
		}
		gosloc.Total += lines
		gosloc.FilePaths[fp] = lines
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
