package queue

type PriorityQueueItems interface {
	Insert(PriorityQueueItem)
	Get() PriorityQueueItem
	Len() int
}
