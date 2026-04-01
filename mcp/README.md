# Hessian MCP Server

MCP (Model Context Protocol) 服务，用于解析 Hessian 序列化数据。

## 项目简介

Hessian 是一种基于 HTTP 的轻量级远程调用协议，由 Caucho 公司开发。本 MCP 服务提供了解析 Hessian V1 序列化数据的能力，可用于：

- 安全研究和漏洞分析
- Java 反序列化漏洞检测
- Hessian 数据结构分析

## 功能特性

- 接收 Base64 编码的 Hessian 序列化数据
- 解析并返回 JSON 格式的结构化结果
- 支持 Hessian V1 协议的各种数据类型：
  - 基本类型：null, boolean, integer, long, double, date
  - 字符串：String
  - 二进制：Binary
  - 集合类型：List, Map
  - 对象引用：Class, Object

## 安装

### 环境要求

- Python >= 3.10
- Go >= 1.21 (用于编译 `hs_serial_parse`)

### 安装步骤

```bash
# 激活虚拟环境（根据你的环境配置）
source /path/to/your/venv/bin/activate

# 安装 MCP 服务
cd mcp
pip install -e .
```

### 编译解析工具

如果需要重新编译 `hs_serial_parse`：

```bash
cd /path/to/hessian_serial_parse
go build -o hs_serial_parse ./cmd
```

## 配置

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `HS_SERIAL_PARSE_PATH` | `hs_serial_parse` 可执行文件路径 | `../hs_serial_parse` |

## 使用方法

### 命令行启动

```bash
# 激活虚拟环境
source /path/to/your/venv/bin/activate

# 使用默认路径
python hessian_mcp.py

# 或指定 hs_serial_parse 路径
HS_SERIAL_PARSE_PATH=/path/to/hs_serial_parse python hessian_mcp.py
```

### Claude Desktop 配置

在 Claude Desktop 配置文件中添加:

**配置文件路径**:
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

**配置内容**:

```json
{
  "mcpServers": {
    "hessian": {
      "command": "/path/to/your/venv/bin/python",
      "args": ["/path/to/hessian_serial_parse/mcp/hessian_mcp.py"],
      "env": {
        "HS_SERIAL_PARSE_PATH": "/path/to/hessian_serial_parse/hs_serial_parse"
      }
    }
  }
}
```

## 工具说明

### parse_hessian

解析 Base64 编码的 Hessian 序列化数据。

**参数**:

| 参数名 | 类型 | 必需 | 说明 |
|--------|------|------|------|
| `data` | string | 是 | Base64 编码的 Hessian 序列化数据 |

**返回值**:

解析后的 JSON 结构，包含以下信息：

- `ClassName`: 对象的 Java 类名
- `Maps`: Map 类型的键值对列表
- `ListType`: List 类型的 Java 类名
- `Args`: List 中的元素列表

**示例**:

```python
import base64

# 读取原始 Hessian 序列化数据
with open("examples/Rome.ser", "rb") as f:
    raw_data = f.read()

# Base64 编码
base64_data = base64.b64encode(raw_data).decode()

# 在 MCP 调用中使用
# result = parse_hessian(data=base64_data)
```

**返回示例**:

```json
{
  "ClassName": "com.example.SomeClass",
  "Maps": [
    {
      "Key": "fieldName",
      "Value": {
        "ClassName": "java.lang.String",
        "Maps": []
      }
    }
  ]
}
```

## 错误处理

| 错误信息 | 原因 | 解决方案 |
|----------|------|----------|
| `hs_serial_parse not found` | 解析工具路径不正确 | 设置 `HS_SERIAL_PARSE_PATH` 环境变量 |
| `Error decoding base64 data` | Base64 编码格式错误 | 检查输入数据是否为有效的 Base64 字符串 |
| `Error parsing Hessian data` | Hessian 数据格式错误 | 检查数据是否为有效的 Hessian 序列化数据 |
| `Parsing timeout` | 解析超时 | 数据可能过大或包含复杂嵌套结构 |

## 测试

```bash
# 激活虚拟环境
source /path/to/your/venv/bin/activate

# 运行测试
cd mcp
python test_mcp.py
```

## 文件结构

```
hessian_serial_parse/
├── cmd/
│   └── main.go           # 命令行工具入口
├── examples/             # 示例文件目录
│   ├── Rome.ser
│   ├── Resin.ser
│   ├── SpringAbstractBeanFactoryPointcutAdvisor.ser
│   ├── SpringPartiallyComparableAdvisorHolder.ser
│   └── XBean.ser
├── utils/
│   ├── encode.go         # JSON 编码工具
│   └── reader.go         # 数据读取工具
├── glb.go                # 常量定义
├── hessian.go            # Hessian 解析核心逻辑
├── hessian_test.go       # Go 测试文件
├── hs_serial_parse       # 编译后的可执行文件
├── go.mod                # Go 模块定义
├── README.md             # 项目说明
└── mcp/
    ├── hessian_mcp.py    # MCP 服务主程序
    ├── pyproject.toml    # Python 项目配置
    ├── README.md         # MCP 服务说明
    └── test_mcp.py       # MCP 测试脚本
```

## 相关链接

- [Hessian Protocol Specification](http://hessian.caucho.com/doc/hessian-serialization.html)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Python SDK](https://github.com/modelcontextprotocol/python-sdk)

## License

MIT License
