package goeureka

import (
	"encoding/json"
	"runtime"
)

func metrics() ([]byte, error) {
	b, err := json.Marshal(refillMetricsMap())
	if err != nil {
		return nil, err
	}
	return b, nil
}

func refillMetricsMap() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return mem
}
