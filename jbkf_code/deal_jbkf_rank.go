package check_jbkf

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/xuri/excelize/v2"
)

func DealRank(path string) {
	// 定义数据结构
	type Record struct {
		IP        string
		GroupRank int // 组内排名
		ID        string
		SID       string
		HID       string
		Day       string
		Score     int
		MaxScore  string
		LSN       string
		FightN    string
		Winn      string
		Losen     string
		Name      string
	}

	// 打开数据文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建记录集合
	ipGroups := make(map[string][]Record)
	currentIP := ""
	dividerPassed := false

	// 正则表达式
	ipRegex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	rowRegex := regexp.MustCompile(`\(\d+\s+rows\)`)

	// 逐行处理文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if line == "" {
			continue
		}

		// 检测IP地址行（数据块开始）
		if ipRegex.MatchString(line) {
			parts := strings.Split(line, " | ")
			if len(parts) > 0 {
				currentIP = parts[0]
			}
			dividerPassed = false
			continue
		}

		// 跳过行数统计行
		if rowRegex.MatchString(line) {
			continue
		}

		// 跳过表头下方的分隔线
		if strings.HasPrefix(line, "----") {
			dividerPassed = true
			continue
		}

		// 处理数据行
		if currentIP != "" && dividerPassed {
			// 清理特殊字符
			cleanLine := strings.Map(func(r rune) rune {
				if unicode.IsGraphic(r) {
					return r
				}
				return -1
			}, line)

			// 分割字段
			fields := strings.Split(cleanLine, "|")
			if len(fields) < 11 {
				continue
			}

			// 清理字段空格
			for i := range fields {
				fields[i] = strings.TrimSpace(fields[i])
			}

			// 转换Score为整型
			score, err := strconv.Atoi(fields[4])
			if err != nil {
				score = 0
			}

			// 创建记录
			record := Record{
				IP:       currentIP,
				ID:       fields[0],
				SID:      fields[1],
				HID:      fields[2],
				Day:      fields[3],
				Score:    score,
				MaxScore: fields[5],
				LSN:      fields[6],
				FightN:   fields[7],
				Winn:     fields[8],
				Losen:    fields[9],
				Name:     strings.Join(fields[10:], " "),
			}

			// 添加到对应的IP组
			ipGroups[currentIP] = append(ipGroups[currentIP], record)
		}
	}

	// 创建所有记录集合
	var allRecords []Record

	// 按IP分组进行排序和排名
	for _, records := range ipGroups {
		// 按Score降序排序
		sort.Slice(records, func(i, j int) bool {
			return records[i].Score > records[j].Score
		})

		// 添加组内排名
		rank := 1
		for i := range records {
			if i > 0 {
				if records[i].Score == records[i-1].Score {
					records[i].GroupRank = records[i-1].GroupRank
				} else {
					records[i].GroupRank = rank
				}
			} else {
				records[i].GroupRank = rank
			}
			rank++
		}

		// 添加到总记录
		allRecords = append(allRecords, records...)
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 设置标题行
	headers := []string{"组内排名", "IP来源", "ID", "SID", "HID", "Day", "Score", "MaxScore", "LSN", "FightN", "Winn", "Losen", "Name"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, h)
	}

	// 设置列宽
	widths := []float64{10, 15, 12, 6, 8, 5, 8, 9, 6, 8, 8, 8, 30}
	for i, width := range widths {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth("Sheet1", col, col, width)
	}

	// 写入数据
	for rowIndex, record := range allRecords {
		rowData := []interface{}{
			record.GroupRank,
			record.IP,
			record.ID,
			record.SID,
			record.HID,
			record.Day,
			record.Score,
			record.MaxScore,
			record.LSN,
			record.FightN,
			record.Winn,
			record.Losen,
			record.Name,
		}

		for colIndex, value := range rowData {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			f.SetCellValue("Sheet1", cell, value)
		}
	}

	// 添加条件格式 - 每组前三名标记
	// top3Style, _ := f.NewStyle(&excelize.Style{
	// 	Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFFF00"}, Pattern: 1},
	// })

	// // 设置条件格式
	// err = f.SetConditionalFormat("Sheet1", "A2:A"+fmt.Sprint(len(allRecords)+1), []excelize.ConditionalFormatOptions{
	// 	{
	// 		Type:     "formula",
	// 		Criteria: "=A2<=3",
	// 		Format:   top3Style,
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println("设置组内前三名条件格式出错:", err)
	// }

	// 保存Excel文件
	if err := f.SaveAs("GroupRanked_Data.xlsx"); err != nil {
		fmt.Println("保存文件失败:", err)
	} else {
		fmt.Printf("成功导出 %d 条数据到 GroupRanked_Data.xlsx\n", len(allRecords))
		fmt.Println("每个IP组内的前三名用黄色标记")
	}
}
