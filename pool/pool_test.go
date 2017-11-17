package pool

import (
	"testing"
	"time"
)

var maxWorkers = 1000

func TestNewPool(t *testing.T) {
	t.Run("pool", func(t *testing.T) {
		pool := NewPool(maxWorkers)

		for i := 1; i <= 100000; i++ {
			j := i
			pool.AddWorker(func(results chan interface{}, next func()) {
				time.Sleep(1 * time.Millisecond)
				next()
				results <- j*j
			})
		}

		pool.Run()

		results := pool.GetResult().([]interface{})
		f := finder{}
		f.set(results)

		for i := 1; i <= 100000; i++ {
			if !f.findValue(i * i) {
				t.Errorf("Invalid result")
			}
		}
	})
}

type finder map[interface{}]bool

func (f *finder) set(values []interface{}) {
	m := make(map[interface{}]bool)
 	for _, value := range values {
		m[value] = true
	}

	*f = m
}

func (f finder) findValue(value interface{}) bool {
	_, ok := f[value]
	return ok
}