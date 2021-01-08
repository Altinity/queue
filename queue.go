package queue

type PriorityQueue interface {
	Insert(item PriorityQueueItem)
	Get() (item PriorityQueueItem, ok bool)
	Done(item PriorityQueueItem)
	Len() int
	Cancel()
}

type priorityQueue struct {
	items         PriorityQueueItems
	waiting       Set
	inProgress    Set
	c             Conditioner
	cancelled     bool
	drainOnCancel bool
}

func New() PriorityQueue {
	return &priorityQueue{
		items:         NewSlicePriorityQueueItems(),
		waiting:       NewMapSet(),
		inProgress:    NewMapSet(),
		c:             NewCond(),
		cancelled:     false,
		drainOnCancel: false,
	}
}

func (q *priorityQueue) Insert(item PriorityQueueItem) {
	q.c.Lock()
	defer q.c.Unlock()

	if q.cancelled {
		// Do not accept items into cancelled queue
		return
	}

	handle := item.Handle()

	if q.waiting.Has(handle) {
		// Do not accept copies
		return
	}

	// Place item as waiting
	q.waiting.Insert(handle)
	if q.inProgress.Has(handle) {
		// In case item is already being processed it's enough to just place it into waiting,
		// it will be prioritised when Done() is called
		return
	}

	// Completely new item, let's prioritize it and signal for waiters to pick it up
	q.items.Insert(item)
	q.c.Signal()
}

func (q *priorityQueue) Get() (item PriorityQueueItem, ok bool) {
	q.c.Lock()
	defer q.c.Unlock()

	for (q.items.Len() == 0) && !q.cancelled {
		// Wait for items or cancellation
		q.c.Wait()
	}

	switch {
	case q.cancelled && q.drainOnCancel:
		if q.items.Len() == 0 {
			return nil, false
		}
	case q.cancelled:
		return nil, false
	}

	item = q.items.Get()
	handle := item.Handle()

	// Move item from waiting to in progress
	q.waiting.Delete(handle)
	q.inProgress.Insert(handle)

	return item, true
}

func (q *priorityQueue) Done(item PriorityQueueItem) {
	q.c.Lock()
	defer q.c.Unlock()

	handle := item.Handle()

	q.inProgress.Delete(handle)

	// In case this item is again waiting for processing (meaning it was re-added during being processed),
	// let's prioritize it and signal for waiters to pick it up
	if q.waiting.Has(handle) {
		q.items.Insert(item)
		q.c.Signal()
	}
}

func (q *priorityQueue) Len() int {
	q.c.Lock()
	defer q.c.Unlock()
	return q.items.Len()
}

func (q *priorityQueue) Cancel() {
	q.c.Lock()
	defer q.c.Unlock()
	q.cancelled = true
	q.c.Broadcast()
}
