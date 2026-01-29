package queue

import (
	"errors"
	"pizzeria/internal/model"
	"sync"
)

var ErrQueueCapacityNotPositive = errors.New("queue capacity not positive")
var ErrQueueClosed = errors.New("queue is closed")

type Queue struct {
	mu     sync.Mutex
	cond   *sync.Cond
	items  []*model.Order
	head   int
	tail   int
	len    int
	cap    int
	closed bool
}

func NewQueue(initCap int) *Queue {
	if initCap <= 0 {
		panic(ErrQueueCapacityNotPositive)
	}
	q := &Queue{
		items: make([]*model.Order, initCap),
		cap:   initCap,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Enqueue(order *model.Order) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	if q.len >= q.cap {
		err := q.resize()
		if err != nil {
			return err
		}
	}

	q.items[q.tail] = order
	q.tail = (q.tail + 1) % q.cap
	q.len++
	q.cond.Broadcast()

	return nil
}

func (q *Queue) Dequeue() (*model.Order, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.len == 0 {
		if q.closed {
			return nil, ErrQueueClosed
		}
		q.cond.Wait()
	}

	item := q.items[q.head]
	q.items[q.head] = nil
	q.head = (q.head + 1) % cap(q.items)
	q.len--
	q.cond.Broadcast()

	return item, nil
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.len
}

func (q *Queue) Cap() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.cap
}

func (q *Queue) Closed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.closed
}

func (q *Queue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	q.closed = true
	q.cond.Broadcast()

	return nil
}

func (q *Queue) resize() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	// TODO Придумать как ресайзить очередь
	newCap := q.len * 2
	newItems := make([]*model.Order, newCap)

	if q.len > 0 {
		if q.head < q.tail {
			copy(newItems, q.items[q.head:q.tail])
		} else {
			n := copy(newItems, q.items[q.head:])
			copy(newItems[n:], q.items[:q.tail])
		}
	}

	q.items = newItems
	q.head = 0
	q.tail = q.len % newCap
	q.cap = newCap
	q.cond.Broadcast()

	return nil
}
