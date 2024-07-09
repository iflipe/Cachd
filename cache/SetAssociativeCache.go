package cache

type SACache struct {
	sets [][2]extraline
	size int
}

type extraline struct {
	tag   uint16
	valid bool
	word  byte
	used  bool
}

func (c *SACache) Init(size int) {
	c.size = size
	c.sets = make([][2]extraline, size/2)
	for i := 0; i < size/2; i++ {
		for j := 0; j < 2; j++ {
			c.sets[i][j] = extraline{0, false, 0, false}
		}
	}
}

func (c *SACache) Size() int {
	return c.size
}

func (c *SACache) Lookup(key uint16) bool {
	kset := key & 0b0000000111111111
	ktag := key &^ kset
	for i := 0; i < 2; i++ {
		if c.sets[kset][i].valid && c.sets[kset][i].tag == ktag {
			return true
		}
	}
	if c.sets[kset][0].used && !c.sets[kset][1].used {
		c.sets[kset][1] = extraline{ktag, true, 0, true}
		c.sets[kset][0].used = false
	} else if !c.sets[kset][0].used {
		c.sets[kset][0] = extraline{ktag, true, 0, true}
		c.sets[kset][1].used = false
	}
	return false
}
