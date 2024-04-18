package cacheMaster

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"log"
	"sync"
)

var (
	OperationCache = NewCache()
	Operations     = map[string]int{"+": 0, "-": 1, "*": 2, "/": 3}
	OperatorByID   = map[int]string{1: "+", 2: "-", 3: "*", 4: "/"}
)

type Cache struct {
	userOperationTimes map[int]map[int]int
	mu                 sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		userOperationTimes: make(map[int]map[int]int),
	}
}

func (c *Cache) Get(userID, operationID int) (int, bool) {

	log.Println("Getting user cache")

	c.mu.Lock()
	defer c.mu.Unlock()

	userTimes, found := c.userOperationTimes[userID]
	if !found {
		return 0, false
	}

	time, found := userTimes[operationID]
	return time, found
}

func (c *Cache) Set(userID, operationID, time int) {

	log.Println("Setting user cache")

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.userOperationTimes[userID]; !ok {
		c.userOperationTimes[userID] = make(map[int]int)
	}

	c.userOperationTimes[userID][operationID] = time
}

func (c *Cache) SetList(userID int, times []int) {

	log.Println("Setting a list of user cache")

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.userOperationTimes[userID]; !ok {
		c.userOperationTimes[userID] = make(map[int]int)
	}

	for operationID, time := range times {
		c.userOperationTimes[userID][operationID] = time
	}
}

func (c *Cache) GetList(userID int) []int {

	log.Println("Getting a list of user cache")

	c.mu.Lock()
	defer c.mu.Unlock()

	operationTimes, found := c.userOperationTimes[userID]
	if !found {
		return []int{}
	}

	var times []int
	for _, time := range operationTimes {
		times = append(times, time)
	}

	log.Println(times)

	return times
}

func LoadOperationTimesIntoCache() error {

	log.Println("Loading a list of user cache")

	userIDs, err := databaseManager.DB.GetUserIDs()
	if err != nil {
		log.Println("Error fetching user IDs from the database:", err)
		return err
	}

	for _, userID := range userIDs {
		times, err := databaseManager.DB.GetTimes(userID)
		if err != nil {
			log.Printf("Error fetching data from the database for userID %d: %v\n", userID, err)
			return err
		}
		OperationCache.SetList(userID, times)
	}

	return nil
}
