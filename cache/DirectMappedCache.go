package cache

type DMCache struct {
	lines []line
	size  int
	wsize int
}

type line struct {
	tag   uint16
	valid bool
	word  []byte
}

func (c *DMCache) Init(size int, ws int) {
	c.size = size
	c.wsize = ws
	c.lines = make([]line, size/ws)
	for i := 0; i < len(c.lines); i++ {
		c.lines[i] = line{0, false, make([]byte, ws)}
	}
}

func (c *DMCache) Size() int {
	return len(c.lines)
}

func (c *DMCache) Lookup(key uint16) bool {
	kblock := (key & (uint16(c.size - 1))) / uint16(c.wsize)
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
