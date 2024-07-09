package cache

type SACache struct {
	sets  [][2]extraline
	size  int
	wsize int
}

type extraline struct {
	tag   uint16
	valid bool
	word  []byte
	used  bool
}

func (c *SACache) Init(size, wordsize int) {
	c.size = size
	c.wsize = wordsize
	c.sets = make([][2]extraline, size/(2*wordsize))
	for i := 0; i < len(c.sets); i++ {
		for j := 0; j < 2; j++ {
			c.sets[i][j] = extraline{0, false, make([]byte, wordsize), false}
		}
	}
}

func (c *SACache) Size() int {
	return c.size
}

func (c *SACache) Lookup(key uint16) bool {
	kset := (key & (uint16(len(c.sets) - 1))) / uint16(c.wsize)
	ktag := key / uint16(c.size)
	for i := 0; i < 2; i++ {
		if c.sets[kset][i].valid && c.sets[kset][i].tag == ktag {
			return true
		}
	}
	if c.sets[kset][0].used && !c.sets[kset][1].used {
		c.sets[kset][1] = extraline{ktag, true, nil, true}
		c.sets[kset][0].used = false
	} else if !c.sets[kset][0].used {
		c.sets[kset][0] = extraline{ktag, true, nil, true}
		c.sets[kset][1].used = false
	}
	return false
}
