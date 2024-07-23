package cache

type SACache struct {
	sets  [][2]SetLine //2-way set associative
	size  int
	wsize int //wordsize in bytes
}

// The struct that represents a line in the cache
type SetLine struct {
	tag   uint16
	valid bool
	word  []byte
	used  bool //used for LRU replacement policy adopted for this experiment
}

// Initializes the instance of DMCache with the appropriate number of zero-valued lines
func (cache *SACache) Init(size, wordsize int) {
	cache.size = size
	cache.wsize = wordsize
	cache.sets = make([][2]SetLine, size/(2*wordsize))
	for i := 0; i < len(cache.sets); i++ {
		for j := 0; j < 2; j++ {
			cache.sets[i][j] = SetLine{0, false, make([]byte, wordsize), false}
		}
	}
}

func (cache *SACache) Size() int {
	return cache.size
}

// Return true if the key is present and valid in the cache or false otherwise
// The key is split into three parts namely the tag, the block and the word
// In this current implementation the 'word' part of the key is not given any treatment since it doesn't
// affect the results. Nonetheless, for consistency, nil value is assigned to the word array every time necessity arises.
func (cache *SACache) Lookup(key uint16) bool {
	//word := key % uint16(c.wsize) //unused variable
	keySet := (key / uint16(cache.wsize)) % uint16(len(cache.sets)) //Set is the middle "(len(cache.sets))" bits
	keyTag := key / uint16(cache.wsize) / uint16(len(cache.sets))   //Tag is the remaining bits
	for i := 0; i < 2; i++ {
		//Tests if there is a hit in any of the two ways of the line
		if cache.sets[keySet][i].valid && cache.sets[keySet][i].tag == keyTag {
			return true
		}
	}
	//Theses lines execute in cases of a miss, either because the block is invalid or the tag is different
	//They are responsible for the implementation of the LRU replacement policy
	if cache.sets[keySet][0].used && !cache.sets[keySet][1].used {
		cache.sets[keySet][1] = SetLine{keyTag, true, nil, true}
		cache.sets[keySet][0].used = false
	} else if !cache.sets[keySet][0].used {
		cache.sets[keySet][0] = SetLine{keyTag, true, nil, true}
		cache.sets[keySet][1].used = false
	}
	return false
}
