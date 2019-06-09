package channels

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-Or(
		sig(3*time.Second),
		sig(2*time.Second),
		sig(1*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
	)

	t.Logf("done after %v", time.Since(start))
}

func TestTee(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := Tee(done, Take(done, Repeat(done, 1, 2), 4))

	for val1 := range out1 {
		t.Logf("out1: %v, out2: %v\n", val1, <-out2)
	}
}

func TestBridge(t *testing.T) {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range Bridge(nil, genVals()) {
		t.Logf("%v ", v)
	}
}
