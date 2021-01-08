package queue

import "context"

type PriorityQueue interface {
	Insert(item PriorityQueueItem)
	Get() (item PriorityQueueItem, ctx context.Context, ok bool)
	Done(item PriorityQueueItem)
	Len() int
	Close()
}

type priorityQueue struct {
	items       PriorityQueueItems
	waiting     Map
	inProgress  Set
	cancelFns   Map
	c           Conditioner
	closed      bool
	drainClosed bool
}

func New() PriorityQueue {
	return &priorityQueue{
		items:       NewSlicePriorityQueueItems(),
		waiting:     NewSimpleMap(),
		inProgress:  NewMapSet(),
		cancelFns:   NewSimpleMap(),
		c:           NewCond(),
		closed:      false,
		drainClosed: false,
	}
}

func (q *priorityQueue) Insert(item PriorityQueueItem) {
	q.c.Lock()
	defer q.c.Unlock()

	if q.closed {
		// Do not accept items into closed queue
		return
	}

	handle := item.Handle()

	// Place item as waiting
	q.waiting.Insert(handle, item)

	if q.inProgress.Has(handle) {
		// In case item is already being processed it's enough to just place it into waiting,
		// it will be prioritised when Done() is called
		fn := q.cancelFns.Get(handle)
		fn.(context.CancelFunc)()
		return
	}

	// Completely new item, let's prioritize it and signal for waiters to pick it up
	q.items.Insert(item)
	q.c.Signal()
}

func (q *priorityQueue) Get() (item PriorityQueueItem, ctx context.Context, ok bool) {
	q.c.Lock()
	defer q.c.Unlock()

	for (q.items.Len() == 0) && !q.closed {
		// Wait for items or being close
		q.c.Wait()
	}

	switch {
	case q.closed && q.drainClosed:
		if q.items.Len() == 0 {
			// Queue drained
			return nil, nil, false
		}
	case q.closed:
		return nil, nil, false
	}

	item = q.items.Get()
	handle := item.Handle()

	// Move item from waiting to in progress
	q.waiting.Delete(handle)
	q.inProgress.Insert(handle)
	c, fn := context.WithCancel(context.Background())
	q.cancelFns.Insert(handle, fn)

	return item, c, true
}

func (q *priorityQueue) Done(item PriorityQueueItem) {
	q.c.Lock()
	defer q.c.Unlock()

	handle := item.Handle()

	q.inProgress.Delete(handle)
	q.cancelFns.Delete(handle)

	// In case this item is again waiting for processing (meaning it was re-added during being processed),
	// let's prioritize it and signal for waiters to pick it up
	if q.waiting.Has(handle) {
		q.items.Insert(q.waiting.Get(handle).(PriorityQueueItem))
		q.c.Signal()
	}
}

func (q *priorityQueue) Len() int {
	q.c.Lock()
	defer q.c.Unlock()
	return q.items.Len()
}

func (q *priorityQueue) Close() {
	q.c.Lock()
	defer q.c.Unlock()
	q.closed = true
	q.c.Broadcast()
}
