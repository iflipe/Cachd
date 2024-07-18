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
	//create an array of the apropriate size to respect the overall size parameter
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
func (c *DMCache) Lookup(key uint16) bool {
	kblock := (key & (uint16(c.size - 1))) / uint16(c.wordSize)
	ktag := key / uint16(c.size)
	tag, valid := c.lines[kblock].tag, c.lines[kblock].valid
	if !valid {
		c.lines[kblock] = line{key / uint16(c.size), true, nil}
		return false
	}
	if tag == ktag {
		return true
	}
	c.lines[kblock] = line{ktag, true, nil}

	return false
}
