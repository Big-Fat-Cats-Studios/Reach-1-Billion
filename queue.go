package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/emirpasic/gods/queues/arrayqueue"
)

type Queue struct {
	q  *arrayqueue.Queue
	mu sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{q: arrayqueue.New()}
}

func (q *Queue) Enqueue(item uint32) { // Changed to uint32 to match your logic
	q.mu.Lock()
	defer q.mu.Unlock()
	q.q.Enqueue(item)
}

func (q *Queue) Dequeue() (uint32, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	val, ok := q.q.Dequeue()
	if !ok {
		return 0, false
	}
	return val.(uint32), true
}

var q = NewQueue()

// Use atomic-friendly types
var globalNumber uint32
var isRunning uint32 // 0 for false, 1 for true

func processQueue() {
	if !atomic.CompareAndSwapUint32(&isRunning, 0, 1) {
		return
	}

	defer atomic.StoreUint32(&isRunning, 0)

	for {
		val, ok := q.Dequeue()
		if !ok {
			break
		}

		// 3. Thread-safe updates to the global number
		switch val {
		case 676767: // reset
			atomic.StoreUint32(&globalNumber, 0)
		case 676766: // multiply 2
			// For complex math, we load, calculate, then store
			old := atomic.LoadUint32(&globalNumber)
			atomic.StoreUint32(&globalNumber, old*2)
		case 767676: // multiply 5
			old := atomic.LoadUint32(&globalNumber)
			atomic.StoreUint32(&globalNumber, old*5)
		case 6767678: // division 2
			old := atomic.LoadUint32(&globalNumber)
			atomic.StoreUint32(&globalNumber, old/2)
		case 7676768: // division 5
			old := atomic.LoadUint32(&globalNumber)
			atomic.StoreUint32(&globalNumber, old/5)
		default:
			atomic.AddUint32(&globalNumber, val)
		}

		fmt.Printf("Processed: %d | Current Total: %d\n", val, atomic.LoadUint32(&globalNumber))
	}
}
