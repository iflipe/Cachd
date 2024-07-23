package cache

/*
# The values related to the memory addresses are represented by the uint16 type
*/

type DMCache struct {
	lines    []line
	size     int //in bytes
	wordSize int //in bytes
}

type line struct {
	tag   uint16
	valid bool
	word  []byte
}

//Initializes the instance of DMCache with the appropriate number of zero-valued lines
func (c *DMCache) Init(size int, ws int) {
	c.size = size
	c.wordSize = ws
	//Creates an array of the apropriate size to respect the overall size parameter
	c.lines = make([]line, size/ws)
	for i := 0; i < len(c.lines); i++ {
		c.lines[i] = line{0, false, make([]byte, ws)}
	}
}

//Returns the size of the cache
func (c *DMCache) Size() int {
	return len(c.lines)
}

//Return true if the key is present and valid in the cache or false otherwise
//The key is split into three parts namely the tag, the block and the word
//In this current implementation the 'word' part of the key is not given any treatment since it doesn't
//affect the results. Nonetheless, for consistency, nil value is assigned to the word array every time necessity arises.
func (cache *DMCache) Lookup(key uint16) bool {
	/*
		The key is split into three parts, the tag, the block and the word
		  tag       block     word
		[______|_____________|____]
		The size of each part is determined by the size of the cache and the wordsize
	*/
	//word := key & (uint16(c.wordSize)-1) //unused variable
	keyBlock := (key & (uint16(cache.size - 1))) / uint16(cache.wordSize)
	keyTag := key / uint16(cache.size)
	tag, valid := cache.lines[keyBlock].tag, cache.lines[keyBlock].valid
	if valid && tag == keyTag {
		return true
	}
	//Theses lines execute in cases of a miss, either because the block is invalid or the tag is different
	cache.lines[keyBlock] = line{keyTag, true, nil}
	return false
}
