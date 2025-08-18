package check_jbkf

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

func DealAwarded(path string) {
	// 打开数据文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 设置表头
	headers := []string{"IP地址", "时间", "标记", "玩家ID", "排名", "时间戳"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	// 设置列宽
	f.SetColWidth("Sheet1", "A", "A", 15) // IP地址
	f.SetColWidth("Sheet1", "B", "B", 20) // 时间
	f.SetColWidth("Sheet1", "C", "C", 20) // 标记
	f.SetColWidth("Sheet1", "D", "D", 15) // 玩家ID
	f.SetColWidth("Sheet1", "E", "E", 8)  // 排名
	f.SetColWidth("Sheet1", "F", "F", 12) // 时间戳

	// 正则表达式匹配数据行
	dataRegex := regexp.MustCompile(`^\./(\d+\.\d+\.\d+\.\d+):(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}),\s+(\w+),\s+(\d+),\s+(\d+),\s+(\d+)$`)

	rowNum := 2 // 从第2行开始写入数据

	// 逐行处理文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 使用正则表达式匹配数据
		matches := dataRegex.FindStringSubmatch(line)
		if len(matches) != 7 { // 整个匹配 + 6个捕获组
			continue
		}

		// 提取匹配结果
		ip := matches[1]
		timeStr := matches[2]
		marker := matches[3]
		playerID := matches[4]
		rank := matches[5]
		timestamp := matches[6]

		// 写入Excel
		data := []string{ip, timeStr, marker, playerID, rank, timestamp}
		for col, value := range data {
			cell, _ := excelize.CoordinatesToCellName(col+1, rowNum)
			f.SetCellValue("Sheet1", cell, value)
		}

		rowNum++
	}

	// 保存Excel文件
	if err := f.SaveAs("GameRankData.xlsx"); err != nil {
		fmt.Println("保存文件失败:", err)
	} else {
		fmt.Printf("成功导出 %d 条数据到 GameRankData.xlsx\n", rowNum-2)
	}
}
