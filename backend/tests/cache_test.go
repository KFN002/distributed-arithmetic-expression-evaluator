package tests

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/cacheMaster"
	"reflect"
	"sync"
	"testing"
)

func TestCache_SetGet(t *testing.T) {
	cache := cacheMaster.NewCache()

	cache.Set(1, 0, 10)
	cache.Set(1, 1, 20)
	cache.Set(2, 0, 30)

	time, found := cache.Get(1, 0)
	if !found {
		t.Errorf("Expected to find time for userID 1 and operationID 0")
	}
	if time != 10 {
		t.Errorf("Expected time to be 10, got %d", time)
	}

	time, found = cache.Get(1, 1)
	if !found {
		t.Errorf("Expected to find time for userID 1 and operationID 1")
	}
	if time != 20 {
		t.Errorf("Expected time to be 20, got %d", time)
	}

	time, found = cache.Get(2, 0)
	if !found {
		t.Errorf("Expected to find time for userID 2 and operationID 0")
	}
	if time != 30 {
		t.Errorf("Expected time to be 30, got %d", time)
	}

	time, found = cache.Get(3, 0)
	if found {
		t.Errorf("Expected not to find time for userID 3 and operationID 0")
	}
}

func TestCache_SetGetList(t *testing.T) {
	cache := cacheMaster.NewCache()

	cache.Set(1, 0, 10)
	cache.Set(1, 1, 20)
	cache.Set(1, 2, 30)
	cache.Set(1, 3, 40)

	cache.SetList(2, []int{1, 2, 3, 4})

	expected := map[int][]int{
		1: {10, 20, 30, 40},
		2: {1, 2, 3, 4},
	}

	for userID, expectedList := range expected {
		result := cache.GetList(userID)
		if !reflect.DeepEqual(result, expectedList) {
			t.Errorf("Expected list for userID %d to be %v, got %v", userID, expectedList, result)
		}
	}
}

func TestCache_SetGetConcurrency(t *testing.T) {
	cache := cacheMaster.NewCache()
	var wg sync.WaitGroup
	numRoutines := 100

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(userID, operationID, time int) {
			defer wg.Done()
			cache.Set(userID, operationID, time)
		}(i, 0, i*10)
	}

	wg.Wait()

	for i := 0; i < numRoutines; i++ {
		time, found := cache.Get(i, 0)
		if !found {
			t.Errorf("Expected to find time for userID %d and operationID 0", i)
		}
		expected := i * 10
		if time != expected {
			t.Errorf("Expected time for userID %d to be %d, got %d", i, expected, time)
		}
	}
}

func TestCache_GetNonExistingUser(t *testing.T) {
	cache := cacheMaster.NewCache()

	time, found := cache.Get(999, 0)
	if found || time != 0 {
		t.Errorf("Expected not to find time for non-existing user, got %d, found: %v", time, found)
	}
}

func TestCache_GetNonExistingOperation(t *testing.T) {
	cache := cacheMaster.NewCache()
	cache.Set(123, 0, 10)

	time, found := cache.Get(123, 1)
	if found || time != 0 {
		t.Errorf("Expected not to find time for non-existing operation, got %d, found: %v", time, found)
	}
}
