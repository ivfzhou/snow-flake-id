# 一、说明

[![codecov](https://codecov.io/gh/ivfzhou/snow-flake-id/graph/badge.svg?token=QYBRAOTH5K)](https://codecov.io/gh/ivfzhou/snow-flake-id)
[![Go Reference](https://pkg.go.dev/badge/gitee.com/ivfzhou/snow-flake-id.svg)](https://pkg.go.dev/gitee.com/ivfzhou/snow-flake-id)

基于 Snowflake 算法的 Go 语言分布式唯一 ID 生成器。

# 二、特性

- **线程安全**：内部使用 `sync.Mutex` 保证并发安全，支持高并发调用
- **时钟回拨处理**：自动检测并处理系统时钟回拨问题，等待时钟追平后继续生成
- **零外部依赖**：仅依赖 Go 标准库，轻量高效
- **高性能**：单节点每毫秒可生成 4096 个不重复 ID（序列号 12 位）

# 三、比特位分布

64 位 ID 的比特位布局如下：

| 位置 | 占用位数 | 说明 | 取值范围 |
|------|----------|------|----------|
| 第 1 位 | 1 bit | 符号位（未使用） | 固定为 0 |
| 第 2~42 位 | 41 bits | 毫秒级时间戳 | 约 69 年可用时间 |
| 第 43~52 位 | 10 bits | 工作机器 ID（machineID） | 0 ~ 1023 |
| 第 53~64 位 | 12 bits | 同一毫秒内的序列号 | 0 ~ 4095 |

```
┌──────────────────────────────────────────────────────────────┐
│ 0 │          41 bits Timestamp         │ 10 bits Machine │ 12 bits Sequence │
└──────────────────────────────────────────────────────────────┘
```

# 四、安装

```bash
go get gitee.com/ivfzhou/snow_flake_id
```

# 五、使用方法

## 5.1 基础用法

```go
package main

import (
    "fmt"
    sfi "gitee.com/ivfzhou/snow_flake_id"
)

func main() {
    // 创建生成器实例，每个节点的 machineID 必须不同（范围：0 ~ 1023）
    generator := sfi.NewGenerator(1)

    // 生成唯一 ID
    id := generator.Generate()
    fmt.Println(id) // 输出：1680307200001234567
}
```

## 5.2 高并发场景

```go
package main

import (
    "fmt"
    "sync"
    sfi "gitee.com/ivfzhou/snow_flake_id"
)

func main() {
    generator := sfi.NewGenerator(1)
    var wg sync.WaitGroup

    // 启动多个 goroutine 并发生成 ID
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            id := generator.Generate()
            fmt.Printf("Generated ID: %d\n", id)
        }()
    }
    wg.Wait()
}
```

# 六、API 参考

### 6.1 `NewGenerator(machineID int64) *Generator`

创建一个 ID 生成器实例。

**参数**：
- `machineID` — 工作机器 ID，取值范围 `0 ~ 1023`。**每个分布式节点的 machineID 必须唯一**。

**异常**：当 `machineID > 1023` 时触发 panic。

### 6.2 `(*Generator) Generate() int64`

生成一个全局唯一的 64 位整型 ID。

**返回值**：生成的雪花算法 ID（`int64` 类型）

**特性**：
- 支持并发调用（内部使用互斥锁保护）
- 当系统时钟发生回拨时，会自动等待时钟恢复后再继续生成
- 当时间戳超出 41 位表示上限时触发 panic

### 6.3 `Generator` 结构体字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `machineID` | `int64` | 工作机器标识 |
| `sequence` | `int64` | 当前毫秒内序列号计数器 |
| `timestamp` | `int64` | 上次生成 ID 时的时间戳（毫秒） |

# 七、注意事项

1. **machineID 唯一性**：在分布式环境中部署时，必须确保每个节点使用不同的 `machineID`（0 ~ 1023），否则可能产生重复 ID。
2. **时钟同步**：请确保服务器时钟保持同步，虽然本实现能容忍短暂的时钟回拨，但频繁或大幅度的时钟回拨会导致性能下降。
3. **时间戳上限**：41 位毫秒时间戳约可支撑 69 年（从初始化时刻算起），超限后会 panic。
4. **QPS 上限**：单节点每毫秒最多生成 4096 个 ID，对应单机 QPS 约为 **409.6 万/秒**。
