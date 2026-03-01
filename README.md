# 一、说明

雪花算法 ID 生成器

64 位比特使用说明：41 位毫秒时间戳，10 位工作机器 ID，12 位序列号，头一位未使用。

# 二、使用

```golang
import sfi "gitee.com/ivfzhou/snow_flake_id"

// 创建一个生成 ID 对象。每个节点的 machineID(ip) 必须不同。
generator := sfi.NewGenerator(machineID)

// 生成唯一id
generator.Generate()

```
