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

func (q *Queue) Enqueue(item int32) { // Changed to uint32 to match your logic
	q.mu.Lock()
	defer q.mu.Unlock()
	q.q.Enqueue(item)
}

func (q *Queue) Dequeue() (int32, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	val, ok := q.q.Dequeue()
	if !ok {
		return 0, false
	}
	return val.(int32), true
}

var q = NewQueue()

// Use atomic-friendly types
var globalNumber uint32
var isRunning uint32
var highScore uint32

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

		switch val {
		case 676767: // reset
			atomic.StoreUint32(&globalNumber, 0)
		case 676766: // multiply 2
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
			// --- SAFE ADDITION / SUBTRACTION ---
			old := atomic.LoadUint32(&globalNumber)
			if val < 0 {
				// Convert val to positive for comparison
				absVal := uint32(-val)
				if absVal >= old {
					atomic.StoreUint32(&globalNumber, 0)
				} else {
					atomic.AddUint32(&globalNumber, uint32(val)) // Two's complement handles subtraction by causing overflow on purpose(numbers wrap around)
				}
			} else {
				atomic.AddUint32(&globalNumber, uint32(val))
			}
		}

		// Update highScore after the math is settled
		currentTotal := atomic.LoadUint32(&globalNumber)
		if currentTotal > atomic.LoadUint32(&highScore) {
			atomic.StoreUint32(&highScore, currentTotal)
		}

		fmt.Printf("Processed: %d | Current Total: %d\n", val, currentTotal)
	}
}
