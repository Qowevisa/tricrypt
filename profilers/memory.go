package profilers

import (
	"os"
	"runtime"
	"runtime/pprof"
)

func GetMemoryProfiler() {
	f, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}
	return
}
