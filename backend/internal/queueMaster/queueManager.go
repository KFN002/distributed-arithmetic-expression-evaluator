package queueMaster

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"sync/atomic"
	"unsafe"
)

var ExpressionsQueue = ExpressionQueue()

// Queue реализация очереди с выражениями через атомики
type Queue interface {
	Enqueue(element models.Expression)
	EnqueueList(data []models.Expression)
	Dequeue() (models.Expression, bool)
}

type QueueNode struct {
	expression models.Expression
	next       unsafe.Pointer
}

type LockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func ExpressionQueue() *LockFreeQueue {
	dummy := &QueueNode{}
	return &LockFreeQueue{
		head: unsafe.Pointer(dummy),
		tail: unsafe.Pointer(dummy),
	}
}

func (q *LockFreeQueue) Enqueue(element models.Expression) {
	newNode := &QueueNode{expression: element}

	for {
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&((*QueueNode)(tail)).next)

		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				if atomic.CompareAndSwapPointer(&((*QueueNode)(tail)).next, nil, unsafe.Pointer(newNode)) {
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
		next := atomic.LoadPointer(&((*QueueNode)(head)).next)

		if head == atomic.LoadPointer(&q.head) {
			if next == nil {
				return models.Expression{}, false
			}
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				return (*QueueNode)(next).expression, true
			}
		}
	}
}
