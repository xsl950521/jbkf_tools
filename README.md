## 结拜跨服赛工具使用说明

```
git clone https://github.com/yourusername/jbkf_tool.git
```

```
cd jbkf_tool
```

```
go build -o jbkf_tool.exe main.go
```

```
jbkf_tool.exe -type=[rank|awarded|pgroup]
```

```
jbkf_tool.exe -help
```

```
jbkf_tool.exe -version
```
jbkf_tool/
├── data/                # 数据目录
│   ├── result_jbkf.txt   # 榜单数据文件
│   └── jbkf_endaward.txt # 已发放数据文件
├── jbkf_code/           # 数据处理代码
│   └── ...              # 数据处理实现
├── main.go              # 主程序入口
└── jbkf_tool.exe        # 可执行文件
