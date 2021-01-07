package queue

import "sort"

type PriorityQueueItem interface {
	Priority() int
}

type PriorityQueueItems interface {
	Insert(PriorityQueueItem)
	Get() PriorityQueueItem
	Len() int
}

type priorityQueueItems struct {
	items []PriorityQueueItem
}

func newPriorityQueueItems() *priorityQueueItems {
	return &priorityQueueItems{}
}

func (i *priorityQueueItems) order() {
	sort.Slice(i.items, func(a, b int) bool {
		return i.items[a].Priority() > i.items[b].Priority()
	})
}

func (i *priorityQueueItems) Insert(item PriorityQueueItem) {
	i.items = append(i.items, item)
	i.order()
}

func (i *priorityQueueItems) Get() (item PriorityQueueItem) {
	item, i.items = i.items[0], i.items[1:]
	i.order()
	return item
}

func (i *priorityQueueItems) Len() int {
	return len(i.items)
}
