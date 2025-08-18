package check_jbkf

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CheckJBKF() {
	outputFile := "missing_jbkf_files.txt"
	rootDir := "jbkf"
	searchStr := "jbkf_rankawardcfend"

	result := checkFiles(rootDir, searchStr)
	writeResults(outputFile, result)
	fmt.Printf("生成报告完成，共发现 %d 个异常文件\n", len(result))
}

func checkFiles(dir, target string) []string {
	var missingFiles []string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if !containsString(path, target) {
				missingFiles = append(missingFiles, path)
			}
		}
		return nil
	})

	return missingFiles
}

func containsString(filename, target string) bool {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), target) {
			return true
		}
	}
	return false
}

func writeResults(output string, data []string) {
	file, _ := os.Create(output)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, path := range data {
		fmt.Fprintln(writer, path)
	}
	writer.Flush()
}
