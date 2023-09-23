/*
 * Copyright (c) 2023 ivfzhou
 * snow-flake-id is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package snow_flake_id

import (
	"sync"
	"time"
)

type generator struct {
	machineID, sequence, timestamp int64
	lock                           sync.Mutex
}

const (
	machineBits   = 10
	sequenceBits  = 12
	timestampBits = 41

	maxMachineID   = 1<<machineBits - 1
	maxSequenceNum = 1<<sequenceBits - 1
	maxTimestamp   = 1<<timestampBits - 1

	timestampShiftBits = machineBits + sequenceBits
	machineIDShiftBits = sequenceBits
)

// NewGenerator 创建一个生成 ID 对象。每个节点的 machineID 必须不同。
func NewGenerator(machineID int64) *generator {
	if machineID > maxMachineID {
		panic("机器 ID 过大")
	}
	return &generator{
		machineID: machineID,
		timestamp: time.Now().UnixMilli(),
	}
}

// Generate 雪花算法生成 ID。41 位毫秒时间戳，10 位工作机器 ID，12 位序列号。
func (g *generator) Generate() int64 {
	g.lock.Lock()
	defer g.lock.Unlock()

	t := time.Now().UnixMilli()
	if t < g.timestamp {
		t = g.nextTime(g.timestamp)
	}

	if t == g.timestamp {
		g.sequence = (g.sequence + 1) & maxSequenceNum
		if g.sequence == 0 {
			g.timestamp = g.nextTime(t + 1)
		}
	} else {
		g.sequence = 0
		g.timestamp = t
	}

	if g.timestamp > maxTimestamp {
		panic("时间超位")
	}

	return g.timestamp<<timestampShiftBits | g.machineID<<machineIDShiftBits | g.sequence
}

func (g *generator) nextTime(timestamp int64) int64 {
	for {
		t := time.Now().UnixMilli()

		if t < timestamp {
			time.Sleep(time.Millisecond)
			continue
		}

		return t
	}
}
