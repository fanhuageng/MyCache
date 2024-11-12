package LRU

import "container/list"

type LRUCache struct {
	maxBytes  int64
	nByets    int64
	ll        *list.List
	cache     map[string]*list.Element
	onEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}
