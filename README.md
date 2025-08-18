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
在工具所在文件夹下打开cmd，根据需要执行以下命令

jbkf_tool.exe -help

jbkf_tool.exe -type=[rank|awarded|pgroup]

jbkf_tool.exe -version
```

### 工具使用注意事项
```
文件名要求:
1、result_jbkf.txt为跨服拉取下来的全服榜单数据
2、jbkf_endaward.txt为GM拉取的已获奖玩家日志；示例：./10.45.136.100:2025-08-08 00:00:00,    jbkf_rankawardcfend,    458661, 15,     807926400
3、result_EU/EST.HKT.txt为GM拉取的各区服在榜group内的玩家数据
所有数据需存储在data文件夹内。
```

### 文件结构
```
jbkf_tool/
├── data/                # 数据目录
│   ├── result_jbkf.txt   # 榜单数据文件
│   └── jbkf_endaward.txt # 已发放数据文件
├── jbkf_code/           # 数据处理代码
│   └── ...              # 数据处理实现
├── main.go              # 主程序入口
└── jbkf_tool.exe        # 可执行文件
```
