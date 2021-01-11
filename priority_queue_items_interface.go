package queue

type Prioritier interface {
	Priority() int
}

type PriorityQueueItem interface {
	Prioritier
	Handle() T
}

type PriorityQueueItems interface {
	// Insert inserts item into the queue
	Insert(PriorityQueueItem)
	Get() PriorityQueueItem
	Len() int
}
