package queue

type Set interface {
	Has(item T) bool
	Insert(item T)
	Delete(item T)
}
