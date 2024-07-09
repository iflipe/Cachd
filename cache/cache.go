package cache

type ICache interface {
	Size() int
	Lookup(key uint16) bool
	Init(size, wordsize int)
}
