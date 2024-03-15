package waiter_test

import (
	"sync"
	"testing"

	"waiter"
)

func wait(w waiter.Waiter) {
	w.Add(10)
	for x := 0; x < 10; x++ {
		go w.Done()
	}
	w.Wait()
}
func Test(t *testing.T) {
	wait(waiter.New())
}
func BenchmarkToy(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wait(waiter.New())
	}
}
func BenchmarkTool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wait(new(sync.WaitGroup))
	}
}
