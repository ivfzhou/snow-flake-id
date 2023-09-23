### 模块说明

雪花算法 ID 生成器

64 位比特使用说明：41 位毫秒时间戳，10 位工作机器 ID，12 位序列号，头一位未使用。

### 快速开始

```golang
import flakeid "github.com/ivfzhou/snow_flake_id"

// 创建一个生成 ID 对象。每个节点的 machineID(ip) 必须不同。
generator := flakeid.NewGenerator(machineID)

// 生成唯一id
generator.Generate()

```

联系电邮：ivfzhou@aliyun.com
