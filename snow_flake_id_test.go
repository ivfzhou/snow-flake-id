package snow_flake_id_test

import (
	"testing"
	"time"

	flakeid "gitee.com/ivfzhou/snow_flake_id"
)

func TestGenerator_Generate(t *testing.T) {
	ch := make(chan int64, 1000)
	m := make(map[int64]bool)
	go func() {
		for i := range ch {
			if m[i] {
				t.Error(i)
			} else {
				m[i] = true
			}
		}
	}()

	g1 := flakeid.NewGenerator(1)
	g2 := flakeid.NewGenerator(2)
	g3 := flakeid.NewGenerator(3)
	g4 := flakeid.NewGenerator(4)
	g5 := flakeid.NewGenerator(5)
	g6 := flakeid.NewGenerator(6)

	for i := 0; i < 1000; i++ {
		switch i % 6 {
		case 0:
			go func() {
				ch <- g2.Generate()
			}()
		case 1:
			go func() {
				ch <- g3.Generate()
			}()
		case 2:
			go func() {
				ch <- g4.Generate()
			}()
		case 3:
			go func() {
				ch <- g5.Generate()
			}()
		case 4:
			go func() {
				ch <- g6.Generate()
			}()
		case 5:
			go func() {
				ch <- g1.Generate()
			}()
		default:
			t.Fatal(i)
		}
	}

	time.Sleep(time.Second * 30)
}

func BenchmarkGenerator_Generate(b *testing.B) {
	g := flakeid.NewGenerator(1)
	for i := 0; i < b.N; i++ {
		g.Generate()
	}
}
