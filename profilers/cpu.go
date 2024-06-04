package profilers

import (
	"os"
	"runtime/pprof"
)

func GetCPUProfiler() func() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	return pprof.StopCPUProfile

}
