package main

import (
    "fmt"
    "sync"
    "github.com/emirpasic/gods/queues/arrayqueue"
)

var q = NewQueue() 
var globalNumber uint32

type Queue struct {
    q  *arrayqueue.Queue
    mu sync.Mutex
}

func NewQueue() *Queue {
    return &Queue{q: arrayqueue.New()}
}

func (q *Queue) Enqueue(item string) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.q.Enqueue(item)
}

func (q *Queue) Dequeue() (interface{}, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()
    return q.q.Dequeue()
}

func processQueue() {
    for {
        val, ok := q.Dequeue()
        if !ok {
            break  
        }
        switch val.(uint32) {
            case 676767:
                globalNumber = 0
            case 676766:
                globalNumber*=2
            case 767676:
                globalNumber*=5
            case 6767678:
                globalNumber/=2
            case 7676768:
                globalNumber/=5
            default:
                globalNumber+=val.(uint32)
        }
        fmt.Println("processing:", val)
    }
}

//676767 - reset
//676766 - multiply 2
//767676 - multiply 5
//6767678 - division 2
//7676768 - division 5