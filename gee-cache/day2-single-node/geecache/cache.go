package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
	// add(key string, value ByteView)		// 返回
	// get(key string) (ByteView, bool)		// 返回 当前的ByteView是什么样子.
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock() // 添加时候上锁.
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil) // 延迟初始化.
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
