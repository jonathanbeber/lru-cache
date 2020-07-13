package lru

import (
	"math/big"
)

// WrappableFunc represents the function signature that `Cache` can wrap.
type WrappableFunc func(int) *big.Int

// Logger represents the loggers `Cache` uses to logging the cache miss, hit
// and eviction operations.
// Its main implementation is the std lib's Logger type.
type Logger interface {
	Printf(format string, v ...interface{})
}

// New returns a Cache instance, given a maximum number of items to keep cached
// (`cacheSize`), the function to be wrapped by the cache (`wrappedFunc`) and a
// Logger (`log`).
func New(cacheSize int, wrappedFunc WrappableFunc, log Logger) Cache {
	return Cache{
		log:         log,
		cacheSize:   cacheSize,
		data:        map[int]*node{},
		wrappedFunc: wrappedFunc,
	}
}

// Cache represents the LRU cache itself. It receives, on its constructor, its
// maximum size, a `WrappableFunc` to keep the cache of and a `Logger`.
//
// Internally it keeps track of the usage order (from most to least recently
// used) of the cached values and holds a data structure for quick access to
// the cached values.
type Cache struct {
	log Logger

	cacheSize   int
	head        *node
	tail        *node
	data        map[int]*node
	wrappedFunc WrappableFunc
}

// Do invokes the wrapped function with a given `key` if the value of `key` is
// not cached. If cached, `Do` returns the cached value.
//
// `Do` is also responsible for updating the internal usage order and evicting
// the least recently used keys when the cache size hits its limit.
func (lc *Cache) Do(key int) *big.Int {
	if v, exists := lc.data[key]; exists {
		lc.log.Printf("cache hit - key: %d\n", key)

		if v.previous != nil {
			v.previous.next = v.next
		}
		if v.next != nil {
			v.next.previous = v.previous
		}

		lc.head.previous = v
		v.next = lc.head

		savedTail := v.previous
		v.previous = nil

		lc.head = v

		if lc.tail == v {
			lc.tail = savedTail
			lc.tail.next = nil
		}
		return v.value
	}

	lc.log.Printf("cache miss - key: %d\n", key)

	v := lc.wrappedFunc(key)

	n := node{
		value:    v,
		key:      key,
		previous: nil,
		next:     lc.head,
	}

	if len(lc.data) == lc.cacheSize {
		lc.log.Printf("cache evicted - key: %d\n", lc.tail.key)
		delete(lc.data, lc.tail.key)
		lc.tail = lc.tail.previous
		lc.tail.next = nil
	}

	if len(lc.data) > 0 {
		lc.head.previous = &n
	}
	lc.head = &n

	if lc.tail == nil {
		lc.tail = &n
	}

	lc.data[key] = &n

	return v
}

type node struct {
	key      int
	value    *big.Int
	previous *node
	next     *node
}
