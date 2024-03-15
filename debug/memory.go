package debug

import (
	"log"
	"runtime"
)

const (
	BYTE_NAME  = "B"
	KBYTE_NAME = "KiB"
	MBYTE_NAME = "MiB"
	GBYTE_NAME = "GiB"
)

const (
	BYTE_TYPE = iota
	KBYTE_TYPE
	MBYTE_TYPE
	GBYTE_TYPE
)

func _getNameFromType(_t uint8) string {
	switch _t {
	case 0:
		return BYTE_NAME
	case 1:
		return KBYTE_NAME
	case 2:
		return MBYTE_NAME
	case 3:
		return GBYTE_NAME
	}
	return "ERROR"
}

type DataShort struct {
	Type       uint8
	Name       string
	Num        uint16
	AfterPoint uint16
	Bytes      uint64
}

func GetDataShort(numOfBytes uint64) DataShort {
	var num, leftBytes, clone, afterPoint uint64
	clone = numOfBytes
	leftBytes = 0
	num = 0
	_type := 0
	for {
		if clone < 1024 {
			break
		}
		num = clone / 1024
		leftBytes = clone % 1024
		clone = num
		_type++
	}
	afterPoint = 0
	afterPointf := (float64(leftBytes)) / 1024.0
	afterPoint = uint64(afterPointf)
	return DataShort{
		Num:        uint16(num),
		AfterPoint: uint16(afterPoint),
		Type:       uint8(_type),
		Name:       _getNameFromType(uint8(_type)),
		Bytes:      numOfBytes,
	}
}

func _printDataShort(name string, d DataShort) {
	log.Printf("%s = %d.%d %s ;; %d\n", name, d.Num, d.AfterPoint, d.Name, d.Bytes)
}

func LogMemUsage() {
	log.Printf("loggin mem Usage")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	aloc := GetDataShort(m.Alloc)
	_printDataShort("Alloc", aloc)
	totalAloc := GetDataShort(m.TotalAlloc)
	_printDataShort("TotalAlloc", totalAloc)
	sys := GetDataShort(m.Sys)
	_printDataShort("Sys", sys)
	log.Printf("\tNumGC = %v\n", m.NumGC)
}
