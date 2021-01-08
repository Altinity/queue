package queue

type PriorityQueueItem interface {
	Priority() int
	Handle() T
}
