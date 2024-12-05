package single_cache

import (
	"math/rand"
	"time"
)

const (
	MaxLevel int     = 16
	P        float32 = 0.5
)

type Node struct {
	key     string
	value   []byte
	forward []*Node
}

type SkipList struct {
	header *Node
	level  int
}

func NewNode(key string, value []byte, level int) *Node {
	return &Node{
		key:     key,
		value:   value,
		forward: make([]*Node, level),
	}
}

func NewSkipList() *SkipList {
	return &SkipList{
		header: NewNode("", nil, 1),
		level:  1,
	}
}

func (sl *SkipList) randomLevel() int {
	level := 1
	rand.Seed(time.Now().UnixNano())
	for rand.Float32() < P && level < MaxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Insert(key string, value []byte) {
	update := make([]*Node, MaxLevel)
	current := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]
	if current != nil && current.key == key {
		current.value = value
		return
	}

	level := sl.randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.header
		}
		sl.level = level
	}

	newNode := NewNode(key, value, level)
	for i := 0; i < level; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}

func (sl *SkipList) Search(key string) ([]byte, bool) {
	current := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
	}
	current = current.forward[0]
	if current != nil && current.key == key {
		return current.value, true
	}
	return nil, false
}
