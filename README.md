# Ledgerd

一个用 Go 编写的极简双重记账命令行工具，通过 JSON 文件实现追加式持久化，便于个人账务入门和后续扩展。

## 特性
- 双重记账：每条凭证必须借贷平衡、至少两条分录、禁止负数。
- 清晰分层：`internal/domain`（领域模型与校验）、`internal/service`（LedgerService 业务规则）、`internal/storage`（文件存储）、`cmd`（Cobra CLI）。
- 文件持久化：所有凭证保存在 `data/journal.json`，遵循 append-only 原则，便于版本管理。

## 项目结构
```
cmd/               // cobra 子命令：add、list、balance
internal/cli/      // App 组装 ledger service
internal/domain/   // JournalEntry / JournalLine 及验证
internal/service/  // LedgerService、余额计算
internal/storage/  // Store 接口与文件实现
data/journal.json  // 默认数据文件（初始化为空数组）
main.go            // 程序入口
go.mod             // 模块定义与依赖
```

## 环境要求
- Go 1.22+
- Git（可选，用于版本管理）

## 构建与安装
```bash
go build ./...
```

构建后会在项目根目录生成可执行文件 `ledgerd`（Windows 下为 `ledgerd.exe`）。可根据需要将其放入 `$GOBIN` 或 PATH 中的任何位置。

## 使用说明

### 全局选项
- `--data <path>`：自定义 JSON 数据文件路径（默认 `data/journal.json`，启动时会自动创建）。

### 添加凭证
可以使用 JSON 文件或命令行参数两种方式：

1. **使用 JSON 文件**
   ```bash
   ./ledgerd add --file entry.json
   ```
   `entry.json` 示例：
   ```json
   {
     "date": "2026-04-05",
     "description": "午餐",
     "lines": [
       {"account": "Expenses:Food", "debit": 50, "credit": 0},
       {"account": "Assets:Cash", "debit": 0, "credit": 50}
     ]
   }
   ```

2. **使用参数输入**
   ```bash
   ./ledgerd add \
     --date 2026-04-05 \
     --description 午餐 \
     --line "Expenses:Food,50,0" \
     --line "Assets:Cash,0,50"
   ```
   `--line` 可重复，多条分录会按顺序写入。

### 列出所有凭证
```bash
./ledgerd list
```
命令会以缩进 JSON 输出当前 `journal.json` 中的所有凭证。

### 查询账户余额
```bash
./ledgerd balance Assets:Cash
```
余额计算遵循账户类别惯例：资产、费用按“借-贷”，负债、收入、权益按“贷-借”，系统通过账户前缀（如 `Assets:`）推断类型。

## 开发/测试建议
- 使用 `go run ./cmd/...` 快速迭代或 `go test ./...`（未来添加测试后）验证逻辑。
- 若需支持更多命令或存储后端，可在 `internal/service` 和 `cmd` 中扩展，保持双重记账的核心验证不变。

