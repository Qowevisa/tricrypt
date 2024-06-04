package shuffle

import "math/rand"

func FastShuffle[T Numeric](ar []T, key []byte) []T {
	if len(key) < 4 {
		return ar
	}
	var seed int64
	seed = int64(key[0])
	seed += int64(key[1]) << 8
	seed += int64(key[2]) << 16
	seed += int64(key[3]) << 24
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(ar), func(i, j int) {
		j = rand.Intn(i + 1)
		ar[i], ar[j] = ar[j], ar[i]
	})
	return ar
}
