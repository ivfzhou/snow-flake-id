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

package snow_flake_id_test

import (
	"sync"
	"testing"

	sfi "gitee.com/ivfzhou/snow_flake_id"
)

func TestGenerator_Generate(t *testing.T) {
	const maximumRoutines = 1000000
	wg := sync.WaitGroup{}
	wg.Add(maximumRoutines)
	generator := sfi.NewGenerator(1)
	set := make(map[int64]struct{})
	ch := make(chan int64, 100)
	for range maximumRoutines {
		go func() {
			defer wg.Done()
			id := generator.Generate()
			ch <- id
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for v := range ch {
		_, has := set[v]
		if has {
			t.Error("id is duplicated", v)
		}
		set[v] = struct{}{}
	}
}
