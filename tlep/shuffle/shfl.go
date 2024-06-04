package shuffle

func allActive(ar []bool) bool {
	for _, b := range ar {
		if b == false {
			return false
		}
	}
	return true
}

func haveTrue(ar []bool) int {
	c := 0
	for _, b := range ar {
		if b {
			c++
		}
	}
	return c
}

func hashIndex(x, i int) int {
	return 13*x + 11*i
}

func zeroValue[T comparable]() T {
	var zero T
	return zero
}

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type ShuffleAlg[T comparable] func([]T, []byte) []T
type ShuffleNumericAlg[T Numeric] func([]T, []byte) []T
type ShuffleNumericAlgWOChanging[T Numeric] func([]T, []byte) []T

func Shuffle[T comparable](msg []T, simmetricKey []byte) []T {
	// L HAS to be the lenght of the bytes, not runes
	l := len(msg)
	var i int
	//
	var ar []T
	i = 0
	for {
		if i == l {
			break
		}
		ar = append(ar, zeroValue[T]())
		i++
	}
	//
	i = 0
	var activIdx []bool
	for {
		if i == l {
			break
		}
		activIdx = append(activIdx, false)
		i++
	}
	//
	c := 0
	sL := len(simmetricKey)
	passCounter := 0
	for {
		if allActive(activIdx) {
			break
		}
		// main part
		var newIdx int
		hashI := 0
		for {
			newIdx = hashIndex(int(simmetricKey[c]), hashI) % l
			if activIdx[newIdx] {
				hashI++
			} else {
				break
			}
		}
		newB := msg[newIdx]
		ar[passCounter] = newB
		activIdx[newIdx] = true
		c++
		if c == sL {
			c = 0
		}
		passCounter++
	}
	return ar
}

func Unshuffle[T comparable](msg []T, simmetricKey []byte) []T {
	// L HAS to be the lenght of the bytes, not runes
	l := len(msg)
	var i int
	//
	var ar []T
	i = 0
	for {
		if i == l {
			break
		}
		ar = append(ar, zeroValue[T]())
		i++
	}
	//
	i = 0
	var activIdx []bool
	for {
		if i == l {
			break
		}
		activIdx = append(activIdx, false)
		i++
	}
	//
	c := 0
	sL := len(simmetricKey)
	passCounter := 0
	for {
		if allActive(activIdx) {
			break
		}
		// main part
		var newIdx int
		hashI := 0
		for {
			newIdx = hashIndex(int(simmetricKey[c]), hashI) % l
			if activIdx[newIdx] {
				hashI++
			} else {
				break
			}
		}
		newB := msg[passCounter]
		ar[newIdx] = newB
		activIdx[newIdx] = true
		c++
		if c == sL {
			c = 0
		}
		passCounter++
	}
	return ar
}
