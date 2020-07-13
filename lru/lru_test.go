package lru

import (
	"fmt"
	"math/big"

	"testing"
)

func TestLruCache(t *testing.T) {
	scenarios := []struct {
		about    string
		key      int
		mustRun  bool
		usedSize int
		tail     int
		log      string
	}{
		{
			about:    "first run - 0 values in cache",
			key:      10,
			mustRun:  true,
			usedSize: 1,
			tail:     10,
			log:      "cache miss - key: 10\n",
			// cache double linked list after run:
			// head: 10 :tail
		},
		{
			about:    "second element - 1 value in cache",
			key:      20,
			mustRun:  true,
			usedSize: 2,
			tail:     10,
			log:      "cache miss - key: 20\n",
			// cache double linked list after run:
			// head: 20|10 :tail
		},
		{
			about:    "third element - 2 values in cache",
			key:      30,
			mustRun:  true,
			usedSize: 3,
			tail:     10,
			log:      "cache miss - key: 30\n",
			// cache double linked list after run:
			// head: 30|20|10 :tail
		},
		{
			about:    "cache hit 10 - 3 values in cache",
			key:      10,
			mustRun:  false,
			usedSize: 3,
			tail:     20,
			log:      "cache hit - key: 10\n",
			// cache double linked list after run:
			// head: 10|30|20 :tail
		},
		{
			about:    "cache hit 20 - 3 values in cache",
			key:      20,
			mustRun:  false,
			usedSize: 3,
			tail:     30,
			log:      "cache hit - key: 20\n",
			// cache double linked list after run:
			// head: 20|10|30 :tail
		},
		{
			about:    "cache hit 10 again - unchanged tail - 3 values in cache",
			key:      10,
			mustRun:  false,
			usedSize: 3,
			tail:     30,
			log:      "cache hit - key: 10\n",
			// cache double linked list after run:
			// head: 10|20|30 :tail
		},
		{
			about:    "one item after cache size - 3 values in cache - drop 30",
			key:      40,
			mustRun:  true,
			usedSize: 3,
			tail:     20,
			log:      "cache miss - key: 40\ncache evicted - key: 30\n",
			// cache double linked list after run:
			// head: 40|10|20 :tail
		},
		{
			about:    "cache hit 20 again - 3 values in cache",
			key:      20,
			mustRun:  false,
			usedSize: 3,
			tail:     10,
			log:      "cache hit - key: 20\n",
			// cache double linked list after run:
			// head: 20|40|10 :tail
		},
		{
			about:    "cache hit 40 again - 3 values in cache",
			key:      40,
			mustRun:  false,
			usedSize: 3,
			tail:     10,
			log:      "cache hit - key: 40\n",
			// cache double linked list after run:
			// head: 40|20|10 :tail
		},
		{
			about:    "cache hit 10 again - 3 values in cache",
			key:      10,
			mustRun:  false,
			usedSize: 3,
			tail:     20,
			log:      "cache hit - key: 10\n",
			// cache double linked list after run:
			// head: 10|40|20 :tail
		},
		{
			about:    "new item after cache - 3 values in cache - drop 20",
			key:      50,
			mustRun:  true,
			usedSize: 3,
			tail:     40,
			log:      "cache miss - key: 50\ncache evicted - key: 20\n",
			// cache double linked list after run:
			// head: 50|10|40 :tail
		},
		{
			about:    "cache 30 miss (comeback) - 3 values in cache - drop 40",
			key:      30,
			mustRun:  true,
			usedSize: 3,
			tail:     10,
			log:      "cache miss - key: 30\ncache evicted - key: 40\n",
			// cache double linked list after run:
			// head: 30|50|10 :tail
		},
	}

	cacheSize := 3
	wf := wrappedFunc{
		ran: false,
	}

	lm := logMock{}

	lruc := New(cacheSize, wf.Do, &lm)
	for _, scenario := range scenarios {

		t.Run(scenario.about, func(t *testing.T) {
			got := lruc.Do(scenario.key)
			// the mocked wrapped function will always return the key * 2
			if got.Cmp(big.NewInt(int64(scenario.key*2))) != 0 {
				t.Fatalf("Expected return to be the same as the key (%d), %d found", scenario.key, got)
			}

			if scenario.mustRun != wf.ran {
				t.Fatalf("Wrong status calling wrapped function, must be %t, got %t", scenario.mustRun, wf.ran)
			}

			if scenario.usedSize != len(lruc.data) {
				t.Fatalf("Expected cache to have size %d after run, got %d", scenario.usedSize, len(lruc.data))
			}

			if scenario.key != lruc.head.key {
				t.Fatalf("Expected cache head to be %d after run, got %d", scenario.key, lruc.head.key)
			}

			if scenario.tail != lruc.tail.key {
				t.Fatalf("Expected cache tail to be %d after run, got %d", scenario.tail, lruc.tail.key)
			}

			if scenario.log != lm.received {
				t.Fatalf("Expected log to be '%s', got '%s'", scenario.log, lm.received)
			}

			// reset mock instances
			lm.received = ""
			wf.ran = false
		})
	}
}

type wrappedFunc struct {
	ran bool
}

func (w *wrappedFunc) Do(i int) *big.Int {
	w.ran = true
	// always return the key * 2
	return big.NewInt(int64(i * 2))
}

type logMock struct {
	received string
}

func (lm *logMock) Printf(format string, v ...interface{}) {
	lm.received += fmt.Sprintf(format, v...)
}
