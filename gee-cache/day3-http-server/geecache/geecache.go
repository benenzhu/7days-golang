package geecache

import (
	"fmt"
	"log"
	"sync"
)

/*
                            是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
*/

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter // 未命中时候回调.
	mainCache cache  // 并发缓存
	// Get(key string) (ByteView,error)
	// 		load(key String) (ByteView, error)
	// 			getLocally(key string) (ByteView, error)
	// 				populateCache(key string, value ByteView)
}

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error) // 通过 key 获取 value
}

// 接口类型函数. GetterFunc 接口型函数的使用场景.
// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g // 上锁, 存在全局的group中.
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key) // 如果没有hit到的话.
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key) // 本地获取 ?
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value) // 在缓存中加入这个东西.
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
