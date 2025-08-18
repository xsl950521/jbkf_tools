package check_jbkf

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func DealPgroup() {
	// 创建Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 创建工作表
	summarySheet := "Summary"
	f.NewSheet(summarySheet)
	f.DeleteSheet("Sheet1") // 删除默认工作表

	// 设置汇总表头
	summaryHeaders := []string{"来源", "主机IP", "ID", "GroupID"}
	for col, header := range summaryHeaders {
		cell := toAlpha(col) + "1"
		f.SetCellValue(summarySheet, cell, header)
	}

	// 设置列宽
	f.SetColWidth(summarySheet, "A", "A", 10) // 来源
	f.SetColWidth(summarySheet, "B", "B", 15) // 主机IP
	f.SetColWidth(summarySheet, "C", "C", 15) // ID
	f.SetColWidth(summarySheet, "D", "D", 15) // GroupID

	// 正则表达式匹配主机行
	hostRegex := regexp.MustCompile(`^([\d.]+)\s+\|\s+SUCCESS\s+\|\s+rc=0\s+>>$`)
	// 正则表达式匹配数据行
	dataRegex := regexp.MustCompile(`^\s*(\d+)\s*\|\s*(\d+)\s*$`)
	// 正则表达式匹配行数统计
	rowsRegex := regexp.MustCompile(`^\(\d+\s+rows\)$`)

	// 处理文件列表
	files := []struct {
		path   string
		source string
	}{
		{"data/result_HKT.txt", "HKT"},
		{"data/result_EST.txt", "EST"},
		{"data/result_EU.txt", "EU"},
	}

	totalRows := 2 // 从第2行开始写入数据

	for _, fileInfo := range files {
		file, err := os.Open(fileInfo.path)
		if err != nil {
			fmt.Printf("无法打开文件 %s: %v\n", fileInfo.path, err)
			continue
		}

		// 为每个来源创建单独的工作表
		sourceSheet := fileInfo.source
		f.NewSheet(sourceSheet)

		// 设置来源工作表表头
		sourceHeaders := []string{"主机IP", "ID", "GroupID"}
		for col, header := range sourceHeaders {
			cell := toAlpha(col) + "1"
			f.SetCellValue(sourceSheet, cell, header)
		}

		// 设置列宽
		f.SetColWidth(sourceSheet, "A", "A", 15) // 主机IP
		f.SetColWidth(sourceSheet, "B", "B", 15) // ID
		f.SetColWidth(sourceSheet, "C", "C", 15) // GroupID

		sourceRowNum := 2 // 来源工作表从第2行开始

		currentHost := ""
		processingData := false

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			// 检查是否是主机行
			if matches := hostRegex.FindStringSubmatch(line); len(matches) > 0 {
				currentHost = matches[1]
				processingData = false
				continue
			}

			// 检查是否是数据开始行
			if strings.Contains(line, "id | groupid") || strings.Contains(line, "----+---------") {
				processingData = true
				continue
			}

			// 检查是否是行数统计
			if rowsRegex.MatchString(line) {
				processingData = false
				continue
			}

			// 处理数据行
			if processingData && currentHost != "" {
				matches := dataRegex.FindStringSubmatch(line)
				if len(matches) == 3 {
					id := matches[1]
					groupID := matches[2]

					// 写入来源工作表
					f.SetCellValue(sourceSheet, "A"+strconv.Itoa(sourceRowNum), currentHost)
					f.SetCellValue(sourceSheet, "B"+strconv.Itoa(sourceRowNum), id)
					f.SetCellValue(sourceSheet, "C"+strconv.Itoa(sourceRowNum), groupID)
					sourceRowNum++

					// 写入汇总工作表
					f.SetCellValue(summarySheet, "A"+strconv.Itoa(totalRows), fileInfo.source)
					f.SetCellValue(summarySheet, "B"+strconv.Itoa(totalRows), currentHost)
					f.SetCellValue(summarySheet, "C"+strconv.Itoa(totalRows), id)
					f.SetCellValue(summarySheet, "D"+strconv.Itoa(totalRows), groupID)
					totalRows++
				}
			}
		}

		file.Close()
		fmt.Printf("已处理文件: %s (%d 行数据)\n", filepath.Base(fileInfo.path), sourceRowNum-2)
	}

	// 设置活动工作表为汇总表
	ss, _ := f.GetSheetIndex(summarySheet)
	f.SetActiveSheet(ss)

	// 保存Excel文件
	if err := f.SaveAs("CombinedGroupData.xlsx"); err != nil {
		fmt.Println("保存文件失败:", err)
	} else {
		fmt.Printf("成功导出 %d 条数据到 CombinedGroupData.xlsx\n", totalRows-2)
	}
}

// 将列索引转换为Excel列名 (0->A, 1->B, ...)
func toAlpha(n int) string {
	if n < 26 {
		return string(rune('A' + n))
	}
	return toAlpha(n/26-1) + string(rune('A'+n%26))
}
