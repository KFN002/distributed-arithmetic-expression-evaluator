package queueMaster

import (
	"distributed-arithmetic-expression-evaluator/backend/models"
	"sync"
)

var ExpressionsQueue = &ConcurrentQueue{}

type Queue interface {
	Enqueue(element models.Expression)
	EnqueueList(data []models.Expression)
	Dequeue() models.Expression
}

type ConcurrentQueue struct {
	queue []models.Expression
	mutex sync.Mutex
}

func (cq *ConcurrentQueue) Enqueue(element models.Expression) {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	cq.queue = append(cq.queue, element)
}

func (cq *ConcurrentQueue) Dequeue() (bool, models.Expression) {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if len(cq.queue) == 0 {
		return false, models.Expression{}
	}
	element := cq.queue[0]
	cq.queue = cq.queue[1:]
	return true, element
}

func (cq *ConcurrentQueue) EnqueueList(data []models.Expression) {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	for _, elem := range data {
		cq.queue = append(cq.queue, elem)
	}
}
