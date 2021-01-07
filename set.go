package queue

type Set interface {
	Has(item t) bool
	Insert(item t)
	Delete(item t)
}

func NewSet() Set {
	return make(MapSet)
}

type empty struct{}
type t interface{}
type MapSet map[t]empty

func (s MapSet) Has(item t) bool {
	_, exists := s[item]
	return exists
}

func (s MapSet) Insert(item t) {
	s[item] = empty{}
}

func (s MapSet) Delete(item t) {
	delete(s, item)
}
