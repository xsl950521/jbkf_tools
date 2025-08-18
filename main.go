package main

import (
	"flag"
	"fmt"
	check_jbkf "jbkf_go/jbkf_code"
	"os"
	"runtime"
	"time"
)

// 程序版本信息
const (
	Version   = "1.1.0"
	BuildDate = "2023-08-18"
)

func main() {
	// 定义命令行参数
	dealType := flag.String("type", "", "指定处理类型: rank/awarded/pgroup")
	help := flag.Bool("help", false, "显示帮助信息")
	version := flag.Bool("version", false, "显示版本信息")
	flag.Parse()

	// 显示版本信息
	if *version {
		fmt.Printf("JBKF数据处理工具 v%s\n", Version)
		fmt.Printf("编译日期: %s\n", BuildDate)
		fmt.Printf("操作系统: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// 显示帮助信息
	if *help || *dealType == "" {
		showHelp()
		os.Exit(0)
	}

	// 显示启动信息
	fmt.Println("========================================")
	fmt.Printf("JBKF数据处理工具 v%s\n", Version)
	fmt.Println("开始处理...")
	fmt.Println("========================================")
	startTime := time.Now()

	// 根据类型执行不同处理
	switch *dealType {
	case "rank":
		fmt.Println("[INFO] 正在处理榜单数据...")
		path := "data/result_jbkf.txt"
		fmt.Printf("[INFO] 使用数据文件: %s\n", path)
		check_jbkf.DealRank(path)
		fmt.Println("[SUCCESS] 榜单数据处理完成")

	case "awarded":
		fmt.Println("[INFO] 正在处理已发放数据...")
		awardPath := "data/jbkf_endaward.txt"
		fmt.Printf("[INFO] 使用数据文件: %s\n", awardPath)
		check_jbkf.DealAwarded(awardPath)
		fmt.Println("[SUCCESS] 已发放数据处理完成")

	case "pgroup":
		fmt.Println("[INFO] 正在处理在榜group的玩家数据...")
		check_jbkf.DealPgroup()
		fmt.Println("[SUCCESS] 在榜group玩家数据处理完成")

	default:
		fmt.Printf("[ERROR] 无效的处理类型 '%s'\n", *dealType)
		fmt.Println("[INFO] 可用选项: rank, awarded, pgroup")
		os.Exit(1)
	}

	// 显示完成信息
	duration := time.Since(startTime)
	fmt.Println("========================================")
	fmt.Printf("[INFO] 处理完成! 耗时: %.2f秒\n", duration.Seconds())
	fmt.Println("========================================")
}

// 显示帮助信息
func showHelp() {
	fmt.Println("JBKF数据处理工具 - 使用说明")
	fmt.Println("========================================")
	fmt.Println("此工具用于处理JBKF游戏数据，支持三种处理模式:")
	fmt.Println()
	fmt.Println("1. 处理榜单数据 (rank):")
	fmt.Println("   分析游戏榜单数据，生成统计报告")
	fmt.Println("   输入文件: data/result_jbkf.txt")
	fmt.Println()
	fmt.Println("2. 处理已发放数据 (awarded):")
	fmt.Println("   分析已发放奖励数据，生成统计报告")
	fmt.Println("   输入文件: data/jbkf_endaward.txt")
	fmt.Println()
	fmt.Println("3. 处理在榜group的玩家数据 (pgroup):")
	fmt.Println("   分析在榜group的玩家数据，生成统计报告")
	fmt.Println()
	fmt.Println("命令行参数:")
	fmt.Println("  -type string   指定处理类型 (必填)")
	fmt.Println("  -help          显示帮助信息")
	fmt.Println("  -version       显示版本信息")
	fmt.Println()
	fmt.Println("使用示例:")
	fmt.Println("  处理榜单数据: jbkf_tool.exe -type=rank")
	fmt.Println("  处理已发放数据: jbkf_tool.exe -type=awarded")
	fmt.Println("  处理在榜group玩家数据: jbkf_tool.exe -type=pgroup")
	fmt.Println()
	fmt.Println("注意: 请确保输入文件位于指定的data目录下")
	fmt.Println("========================================")
}
