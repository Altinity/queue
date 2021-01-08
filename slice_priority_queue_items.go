package queue

import "sort"

type SlicePriorityQueueItems struct {
	items []PriorityQueueItem
}

func NewSlicePriorityQueueItems() *SlicePriorityQueueItems {
	return &SlicePriorityQueueItems{}
}

func (i *SlicePriorityQueueItems) order() {
	sort.Slice(i.items, func(a, b int) bool {
		return i.items[a].Priority() > i.items[b].Priority()
	})
}

func (i *SlicePriorityQueueItems) Insert(item PriorityQueueItem) {
	i.items = append(i.items, item)
	i.order()
}

func (i *SlicePriorityQueueItems) Get() (item PriorityQueueItem) {
	item, i.items = i.items[0], i.items[1:]
	i.order()
	return item
}

func (i *SlicePriorityQueueItems) Len() int {
	return len(i.items)
}
