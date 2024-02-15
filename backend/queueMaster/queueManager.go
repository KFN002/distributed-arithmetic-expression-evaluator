package queueMaster

import (
	"distributed-arithmetic-expression-evaluator/backend/models"
	"sync/atomic"
	"unsafe"
)

var ExpressionsQueue = NewLockFreeQueue()

type Queue interface {
	Enqueue(element models.Expression)
	EnqueueList(data []models.Expression)
	Dequeue() (models.Expression, bool)
}

type Node struct {
	expression models.Expression
	next       unsafe.Pointer
}

type LockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewLockFreeQueue() *LockFreeQueue {
	dummy := &Node{}
	return &LockFreeQueue{
		head: unsafe.Pointer(dummy),
		tail: unsafe.Pointer(dummy),
	}
}

func (q *LockFreeQueue) Enqueue(element models.Expression) {
	newNode := &Node{expression: element}

	for {
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&((*Node)(tail)).next)

		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				if atomic.CompareAndSwapPointer(&((*Node)(tail)).next, nil, unsafe.Pointer(newNode)) {
					atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(newNode))
					return
				}
			} else {
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		}
	}
}

func (q *LockFreeQueue) EnqueueList(data []models.Expression) {
	for _, expr := range data {
		q.Enqueue(expr)
	}
}

func (q *LockFreeQueue) Dequeue() (models.Expression, bool) {
	for {
		head := atomic.LoadPointer(&q.head)
		next := atomic.LoadPointer(&((*Node)(head)).next)

		if head == atomic.LoadPointer(&q.head) {
			if next == nil {
				return models.Expression{}, false
			}
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				return (*Node)(next).expression, true
			}
		}
	}
}
