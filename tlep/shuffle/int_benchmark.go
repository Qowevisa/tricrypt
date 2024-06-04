package shuffle

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

const (
	n          = 5
	samplesNum = 300000
)

func factorial(n int) *big.Int {
	return factorial_body(n, big.NewInt(1))
}

func factorial_body(n int, acc *big.Int) *big.Int {
	if n == 0 {
		return acc
	}
	acc.Mul(acc, big.NewInt(int64(n)))
	return factorial_body(n-1, acc)
}

func ChiSquare(obs, exp []float64) (float64, error) {
	if len(obs) != len(exp) {
		return 0, BENCHMARK_LEN_DN_MATCH
	}
	var result float64
	result = 0
	for i, a := range obs {
		b := exp[i]
		if a == 0 && b == 0 {
			continue
		}
		result += (a - b) * (a - b) / b
	}
	return result, nil
}

func bigChiSquare(observed []int, expected *big.Float) *big.Float {
	chiSquare := big.NewFloat(0)
	for _, obs := range observed {
		obsFloat := big.NewFloat(float64(obs))
		diff := new(big.Float).Sub(obsFloat, expected)
		squaredDiff := new(big.Float).Mul(diff, diff)
		term := new(big.Float).Quo(squaredDiff, expected)
		chiSquare.Add(chiSquare, term)
	}
	return chiSquare
}

func GetBenchmarkN() int {
	return n
}

func GetBenchmarkForShuffle[T Numeric](alg ShuffleNumericAlg[T], keyLen int, fullDebug bool) error {
	perms := map[[n]T]int{}
	for i := 0; i < samplesNum; i++ {
		if fullDebug {
			log.Printf("Running %d sample\n", i)
		}
		arr := [n]T{zeroValue[T]()}
		for i := 0; i < n; i++ {
			arr[i] = T(i)
		}
		//
		dummyKey := make([]byte, keyLen)
		keyN, err := rand.Read(dummyKey)
		if err != nil {
			return err
		}
		if keyN != len(dummyKey) {
			return BENCHMARK_INTERNAL_ERROR
		}
		alg(arr[:], dummyKey)
		//
		perms[arr]++
	}
	fact := factorial(n)
	factFloat := new(big.Float).SetInt(fact)

	// Calculate expected value
	samplesFloat := big.NewFloat(float64(samplesNum))
	expectedBig := new(big.Float).Quo(samplesFloat, factFloat)
	if fullDebug {
		log.Printf("expectedBig is %f\n", expectedBig)
		log.Printf("factorial is %d\n", factorial(n))
	}

	observed := make([]int, 0, len(perms))
	for _, count := range perms {
		// log.Printf("Getting count from perm = %d\n", count)
		observed = append(observed, count)
	}

	val := bigChiSquare(observed, expectedBig)
	fmt.Printf("Chi-Square value: %s\n", val.Text('f', 10))
	return nil
}
