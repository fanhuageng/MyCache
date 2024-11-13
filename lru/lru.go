package LRU

import "container/list"

type LRUCache struct {
	maxBytes  int64
	nByets    int64
	ll        *list.List
	cache     map[string]*list.Element
	onEvicted func(key string, value Value) // 删除缓存的回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int // 返回值所占的内存大小
}

func New(maxBytes int64, onEvicted func(string, Value)) *LRUCache {
	return &LRUCache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (c *LRUCache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront((ele))
		kv := ele.Value.(*entry) // 将entry类型转换成*entry
		return kv.value, ok
	}
	return
}

func (c *LRUCache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nByets -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *LRUCache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nByets += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nByets += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nByets {
		c.RemoveOldest()
	}
}

func (c *LRUCache) Len() int {
	return c.ll.Len()
}
