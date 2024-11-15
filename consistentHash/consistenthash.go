package consistentHash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash            Hash
	virtualMultiple int            // 虚拟节点倍数
	hashLink        []int          // 哈希环
	hashMap         map[int]string // 虚拟节点与真实节点的映射
}

func New(virtualMultiple int, fn Hash) *Map {
	m := &Map{
		virtualMultiple: virtualMultiple,
		hash:            fn,
		hashMap:         make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(nodes ...string) {
	// 遍历真实节点
	for _, node := range nodes {
		for i := 0; i < m.virtualMultiple; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + node))) // 计算虚拟节点的哈希值
			m.hashLink = append(m.hashLink, hash)               // 虚拟节点添加到哈希环上
			m.hashMap[hash] = node                              // 虚拟节点映射真实节点
		}
	}
	sort.Ints(m.hashLink)
}

func (m *Map) Get(key string) string {
	if len(m.hashLink) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key))) // 计算key的哈希值
	// 顺时针找第一个匹配的虚拟节点的下标
	idx := sort.Search(len(m.hashLink), func(i int) bool {
		return m.hashLink[i] >= hash
	})
	return m.hashMap[m.hashLink[idx%len(m.hashLink)]] // 映射得到真实的节点
}
