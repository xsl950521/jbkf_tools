## 结拜跨服赛工具使用说明

### 获取工程代码

```
git clone https://github.com/xsl950521/jbkf_tool.git
```

### 编辑代码
```
cd jbkf_tool

go build -o jbkf_tool.exe main.go
```


### 工具使用
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
>├── data/                # 数据目录
>│   >>├── result_jbkf.txt   # 榜单数据文件
>│   >>└── jbkf_endaward.txt # 已发放数据文件
├── jbkf_code/           # 数据处理代码
│   └── ...              # 数据处理实现
├── main.go              # 主程序入口
└── jbkf_tool.exe        # 可执行文件
