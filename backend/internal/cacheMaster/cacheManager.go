package cacheMaster

import "sync"

var (
	OperationCache = NewCache()
	Operations     = map[string]int{"+": 0, "-": 1, "*": 2, "/": 3}
)

// Cache Кэш с данными операций, чтобы каждый раз не лезть в базу данных
type Cache struct {
	operationTimes map[int]int
	mu             sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		operationTimes: make(map[int]int),
	}
}

func (c *Cache) Get(operationId int) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	time, found := c.operationTimes[operationId]
	return time, found
}

func (c *Cache) Set(operationID int, time int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.operationTimes[operationID] = time
}

func (c *Cache) SetList(times []int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for operationID, time := range times {
		c.operationTimes[operationID] = time
	}
}
