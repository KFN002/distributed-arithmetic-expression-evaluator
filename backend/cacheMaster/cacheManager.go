package cacheMaster

var (
	OperationCache = NewCache()
	Operations     = map[string]int{"+": 1, "-": 2, "*": 3, "/": 4}
)

type Cache struct {
	operationTimes map[int]int
}

func NewCache() *Cache {
	return &Cache{
		operationTimes: make(map[int]int),
	}
}

func (c *Cache) Get(operationId int) (int, bool) {
	time, found := c.operationTimes[operationId]
	return time, found
}

func (c *Cache) Set(operationID int, time int) {
	c.operationTimes[operationID] = time
}

func (c *Cache) SetList(times []int) {
	for operationID, time := range times {
		c.operationTimes[operationID] = time
	}
}
