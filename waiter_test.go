package waiter_test

import (
	"sync"
	"testing"

	"waiter"

	. "github.com/mdwhatcott/funcy"
	"github.com/mdwhatcott/go-set/v2/set"
	"github.com/mdwhatcott/testing/should"
)

func wait(factory func() waiter.Waiter) chan int {
	result := make(chan int, 100)
	defer close(result)

	sut := factory()
	sut.Add(10)
	for x := 0; x < 10; x++ {
		go func(x int) { defer sut.Done(); result <- x }(x)
	}

	helper := factory()
	helper.Add(10)
	for x := 0; x < 10; x++ {
		go func(x int) { defer helper.Done(); sut.Wait(); result <- 10 + x }(x)
	}
	helper.Wait()

	return result
}
func Test(t *testing.T) {
	result := Drain(wait(waiter.New))
	should.So(t, len(result), should.Equal, 20)
	should.So(t, set.Of(Take(10, result)...), should.Equal, set.Of(Range(0, 10)...))
	should.So(t, set.Of(Drop(10, result)...), should.Equal, set.Of(Range(10, 20)...))
	should.So(t, SortAscending(ByNumericValue[int], result), should.Equal, Range(0, 20))
}

func BenchmarkToy(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wait(waiter.New)
	}
}
func BenchmarkTool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wait(func() waiter.Waiter { return new(sync.WaitGroup) })
	}
}

/*
goos: darwin
goarch: arm64
pkg: waiter
BenchmarkToy
BenchmarkToy-10     	     984	   1473062 ns/op	    2244 B/op	      47 allocs/op
BenchmarkTool
BenchmarkTool-10    	  179563	      6375 ns/op	    2209 B/op	      43 allocs/op
PASS
*/
