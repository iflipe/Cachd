package cache

//Interface to make it easier to test both types of cache
type ICache interface {
	Size() int
	Lookup(key uint16) bool
	Init(size, wordsize int)
}
