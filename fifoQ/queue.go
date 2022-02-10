package fifoQ

import (
	"fmt"
	"sync"
)

const defaultSize = 32

type Queue struct {
	rw         sync.RWMutex
	head, rear int
	buf        []interface{}
	empty      bool // when head == rear, this bit is to indicate whether the queue is empty or full
}

func New(size int) *Queue {
	q := &Queue{}
	if size == 0 {
		size = defaultSize
	}
	if size < 0 {
		panic(fmt.Sprintf("negative size: pass 0 to use default size (%v)", defaultSize))
	}
	q.buf = make([]interface{}, size, size)
	q.head, q.rear = 0, 0
	q.empty = true
	return q
}

func (q *Queue) Size() int {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return len(q.buf)
}

func (q *Queue) Cap() int {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return cap(q.buf)
}

func (q *Queue) Enqueue(e interface{}) bool {
	q.rw.Lock()
	defer q.rw.Unlock()
	if (q.head == q.rear) && (q.empty == false) {
		// the queue is full
		return false
	}
	q.buf[q.rear] = e
	q.rear = (q.rear + 1) % cap(q.buf)
	if q.head == q.rear {
		q.empty = false
	}
	return true
}

func (q *Queue) Dequeue() (interface{}, bool) {
	q.rw.Lock()
	defer q.rw.Unlock()
	if q.head == q.rear && q.empty {
		// the queue is empty
		return nil, false
	}
	v := q.buf[q.head]
	q.head = (q.head + 1) % cap(q.buf)
	if q.head == q.rear {
		q.empty = true
	}
	return v, true
}

func (q *Queue) Peek() (interface{}, bool) {
	q.rw.RLock()
	defer q.rw.RUnlock()
	if q.head == q.rear && q.empty {
		return nil, false
	}
	return q.buf[q.head], true
}

func (q *Queue) Full() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()
	if (q.head == q.rear) && (q.empty == false) {
		return true
	}
	return false
}

func (q *Queue) Empty() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()
	if (q.head == q.rear) && q.empty {
		return true
	}
	return false
}
