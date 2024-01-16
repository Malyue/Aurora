package loadbalance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashBalance struct {
	mu       sync.RWMutex
	hash     Hash
	replicas int
	keys     UInt32Slice
	hashMap  map[uint32]string

	conf LoadBalanceConf
}

func (c *ConsistentHashBalance) NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

func (c *ConsistentHashBalance) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param at least 1")
	}
	addr := params[0]
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr
	}

	sort.Sort(c.keys)
	return nil
}

func (c *ConsistentHashBalance) Get(key string) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("node is empty")
	}
	hash := c.hash([]byte(key))

	// find the first index of the node which hash >= hash(key)
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})

	if idx == len(c.keys) {
		idx = 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hashMap[c.keys[idx]], nil
}

func (c *ConsistentHashBalance) SetConf(conf LoadBalanceConf) {
	c.conf = conf
}

func (c *ConsistentHashBalance) Update() {
	if conf, ok := c.conf.(*CheckConf); ok {
		c.keys = nil
		c.hashMap = map[uint32]string{}
		for _, ip := range conf.GetConf() {
			c.Add(strings.Split(ip, ",")...)
		}
	}
}
