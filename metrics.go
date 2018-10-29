package goeureka

import (
	"runtime"
)

func metrics() runtime.MemStats {

	return refillMetricsMap()
}

func refillMetricsMap() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return mem
}