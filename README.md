# CSV RSA-PSS 签名工具

这是一个用于CSV文件RSA-PSS签名的命令行工具，支持Java生成的PKCS#8格式私钥。

## 功能特点

- 🔹 逐行读取CSV文件，对第一列数据进行RSA-PSS签名
- 🔹 支持Java生成的PKCS#8格式私钥
- 🔹 并发处理，充分利用多核CPU
- 🔹 实时显示执行进度
- 🔹 流式处理，支持大规模数据（已测试40万行）
- 🔹 低内存占用，适合处理百万级数据

## 技术栈

- **语言**：Go 1.20+
- **算法**：RSA-PSS + SHA256
- **依赖**：标准库，无第三方依赖

## 使用方法

### 编译程序

```bash
go build -o signcsv.exe
```

### 运行程序

```bash
# 基本用法
.\signcsv.exe -i input.csv -o output.csv -k private.pem

# 指定并发数
.\signcsv.exe -i input.csv -o output.csv -k private.pem -c 8
```

### 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-i` | 输入CSV文件路径 | `test.csv` |
| `-o` | 输出CSV文件路径 | `output.csv` |
| `-k` | RSA私钥文件路径 | `private.pem` |
| `-c` | 并发处理的goroutine数量 | `4` |

## 私钥格式要求

工具只支持Java生成的**PKCS#8格式**私钥，私钥文件应满足：

```
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQ...
-----END PRIVATE KEY-----
```

### 生成私钥（Java）

```bash
# 使用Java keytool生成密钥库
keytool -genkeypair -alias test -keyalg RSA -keysize 2048 -storetype PKCS12 -keystore keystore.p12 -validity 3650

# 使用OpenSSL导出PKCS#8私钥
openssl pkcs12 -in keystore.p12 -nodes -nocerts -out private.pem
```

## 模板说明

当前版本使用固定模板：

```
trade_no=%s&version=1.0
```

- `%s` 会被CSV第一列数据替换
- 生成的待签名内容格式：`trade_no=123456789012&version=1.0`

## 输出格式

输出CSV文件包含两列：

| 列名 | 说明 |
|------|------|
| `out_trade_no-String` | 原始第一列数据 |
| `sign-String` | RSA-PSS签名结果（Base64编码） |

## 性能表现

| 数据规模 | 并发数 | 处理时间 | 处理速度 |
|----------|--------|----------|----------|
| 10,000行 | 4 | ~1.5秒 | ~6,600行/秒 |
| 100,000行 | 8 | ~13秒 | ~7,700行/秒 |
| 400,000行 | 8 | ~55秒 | ~7,300行/秒 |

## 项目结构

```
query/
├── csv.go          # CSV处理逻辑
├── main.go         # 主程序入口
├── rsa.go          # RSA签名逻辑
├── README.md       # 项目说明文档
├── go.mod          # Go模块配置
└── text.csv        # 示例输入文件
```

## 开发说明

### 核心函数

- `ProcessCSVStream()`：流式处理CSV文件
- `LoadPrivateKey()`：加载RSA私钥
- `RSA_PSS_Sign()`：执行RSA-PSS签名

### 扩展建议

1. 支持自定义模板
2. 支持多种哈希算法
3. 支持公钥验证签名
4. 支持更多私钥格式

## 注意事项

1. 🔐 私钥文件包含敏感信息，请勿上传到GitHub
2. ⚠️ 程序仅处理CSV第一列数据，其他列保持不变
3. 📁 输入文件需包含表头行
4. 💡 建议使用8-12个并发数，充分利用多核CPU

## 示例

### 输入文件（text.csv）

```csv
out_trade_no-String
654654654654
564646546464
464654564654
```

### 输出文件（output.csv）

```csv
out_trade_no-String,sign-String
654654654654,f3uXnPzWTOlnMR0zfHrA+DzbbfgcmIvzxmMDIQZxHNk+8jCpCSjC2HU6C4ZaX6e+yfSlq57267PxkgaJVH6FLi+zY0SDzkRw1+PhFALqu1j20c5w6T+mXuQke8QcjRiKABuQ8Df6ySygfoScC3XBC+erpH+VuOiES5+Ih6Tao9TTvMjNQPQWj89c5r2gHUz/8aFQegRAln8X0n2pn4Rx4CgkMEGkeKbOm9uzhe7afXWxk5nWWwccu1Ari4UugIyOme12ewxJ9cJKn2Dd3fplN6T17qKWC4JHGqOugg7avhPntULfRFym2X6f23EwWR93Ylpn6cCNEZTWQFpbgXBwng==
564646546464,UPEXfdzcDwSa5EnOtYXRr+BP2VQPIGjlG2c/iDsoVeNFCa+lcoOQxzRGX05u4LD4FgBiBxCc3ACL5t7m3DnfVytmf2efrTFBGi18FxBrvdM36OtW2GxUpHenfd0KFMJbAB2erNDC9ihOjboPSBsXXeqmcDgAI78UTcvLV7teFhlDodDtfgCxGNwDNm6s2GOYc7BSNI/kiuc22hAEGoYXFKlEjo92EpNSyWD/vtZ+afeNaDfzVoZz5bPY6Zhp8IHIuyD7QUXuLcjk8eqg/pUmRj9a4RWqdN7NOxtfbgHSXZpEjeZ+2Is2LZlxZXR+O+/qGEhxEDwD3+0gMKAgLpEZYg==
464654564654,Ag/qBPCIV5ABfcHutOMY41rCVGYe9rV/Uk5CYdVbjynE76RZn3cPkXO0qxjjOTDVV4xDTIPuuiT+cQCV+gPMbD6qfBcMCPMtCSUVWCZMX50tx6FrVssNMTGFMRZ9j1xOCw6C+uZ//rLRcjo0XSRUngiriYbFF/OYZcpX6M4dj9EsQNRz24lQNCFa3kutg0OGfSDzRvEYWyCVp5Eg16g0il7g87Iz9jCyYvkvCXSfyZ3+Eoq9X1JL2E8ucuynQngkNEgnm6QnBKP0zp56Ac2nGfLIqfFjKHm8Y86iSwK93kxo19FCNTrf7gLle3gbEGgsFoWh1ax0NDhTr0MSCmSeyQ==
```

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！
