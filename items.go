package queue

type PriorityQueueItems interface {
	// Insert inserts item into the queue
	Insert(PriorityQueueItem)
	Get() PriorityQueueItem
	Len() int
}
