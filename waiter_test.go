package waiter_test

import (
	"sync"
	"testing"

	"waiter"

	. "github.com/mdwhatcott/funcy"
	"github.com/mdwhatcott/go-set/v2/set"
	"github.com/mdwhatcott/testing/should"
)

func wait() chan int {
	result := make(chan int, 100)
	defer close(result)

	sut := waiter.New()
	sut.Add(10)
	for x := 0; x < 10; x++ {
		go func(x int) { defer sut.Done(); result <- x }(x)
	}

	var helper sync.WaitGroup
	helper.Add(10)
	for x := 0; x < 10; x++ {
		go func(x int) { defer helper.Done(); sut.Wait(); result <- 10 + x }(x)
	}
	helper.Wait()

	return result
}
func Test(t *testing.T) {
	result := Drain(wait())
	should.So(t, len(result), should.Equal, 20)
	should.So(t, set.Of(Take(10, result)...), should.Equal, set.Of(Range(0, 10)...))
	should.So(t, set.Of(Drop(10, result)...), should.Equal, set.Of(Range(10, 20)...))
	should.So(t, SortAscending(ByNumericValue[int], result), should.Equal, Range(0, 20))
}
