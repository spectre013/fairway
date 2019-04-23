package fairway

import (
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
	"os"
	"runtime"
	"strings"
	"time"
)

var names = []string{
	"jvm.memory.max",
	"jvm.threads.states",
	"jvm.gc.memory.promoted",
	"system.load.average.1m",
	"http.server.requests",
	"jvm.memory.used",
	"jvm.gc.max.data.size",
	"jvm.gc.pause",
	"jvm.memory.committed",
	"system.cpu.count",
	"jvm.buffer.memory.used",
	"jvm.threads.daemon",
	"system.cpu.usage",
	"jvm.gc.memory.allocated",
	"jvm.threads.live",
	"jvm.threads.peak",
	"process.uptime",
	"process.cpu.usage",
	"jvm.gc.live.data.size",
	"jvm.buffer.count",
	"jvm.buffer.total.capacity",
	"process.start.time",
}

var methods = map[string]func() []Measurement{
	"process.uptime":            uptime,
	"jvm.threads.states":        processsThreads,
	"process.start.time":        sinceStart,
	"jvm.memory.used":           memoryUsed,
	"jvm.gc.max.data.size":      jvmGC,
	"jvm.gc.memory.promoted":    jvmGC,
	"jvm.gc.live.data.size":     jvmGC,
	"jvm.gc.pause":              jvmGCPause,
	"jvm.memory.committed":       memoryComitted,
	"system.cpu.count":          cpuCount,
	"jvm.buffer.memory.used":    buffer,
	"jvm.threads.daemon":        threads,
	"system.cpu.usage":          cpuUsage,
	"jvm.buffer.count":          buffer,
	"jvm.buffer.total.capacity": buffer,
	"jvm.threads.live":          processsThreads,
	"process.cpu.usage":         processsCPU,
	"jvm.memory.max":			 memoryMax,
	"jvm.threads.peak":			 processsThreads,
}

var mem runtime.MemStats
var MAXGC uint64
var p process.Process

type Metric struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	BaseUnit      string        `json:"baseUnit"`
	Measurement   []Measurement `json:"measurements"`
	AvailableTags []Tag         `json:"availableTags"`
}

type Measurement struct {
	Statistic string      `json:"statistic"`
	Value     interface{} `json:"value"`
}

type Tag struct {
	Tag    string   `json:"tag"`
	Values []string `json:"values"`
}

func metrics(metric string,query map[string][]string) ([]byte, error) {
	MAXGC = 0
	runtime.ReadMemStats(&mem)
	p = process.Process{Pid:int32(os.Getpid())}
	var m Metric
	if metric != "" {
		if function, ok := methods[metric]; ok {
			logger.Info("Running: ", metric, "()")
			m = createMetric(metric, metricValues[metric])
			m.Measurement = runMetric(function)
			m.AvailableTags = getTags(query)

		} else {
			return nil, errors.New("method not found")
		}
	} else {
		return toJson(map[string][]string{"names": names}), nil
	}

	return toJson(m), nil
}

func getTags(query map[string][]string) []Tag {
	t := make([]Tag, 0)
	for _,a := range query["tags"] {
		tgs := strings.Split(a,":")
		tag := Tag{Tag:tgs[0],Values:[]string{tgs[1]}}
		t = append(t, tag)
	}
	return t
}

func threads() []Measurement {
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": 0.0})
	return createMeaturement(ms)
}

func uptime() []Measurement {
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": time.Since(startTime).Seconds()})
	return createMeaturement(ms)
}

func sinceStart() []Measurement {
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": startTime})
	return createMeaturement(ms)
}

func memoryUsed() []Measurement {
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": mem.HeapInuse})
	return createMeaturement(ms)
}

func jvmGC() []Measurement {
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "COUNT", "stat": 0.0})
	return createMeaturement(ms)
}

func jvmGCPause() []Measurement {
	//PauseTotalNs
	pause := int64(mem.PauseTotalNs)
	//Don't deplay GC Pause until it surpasses 100000 nanoseconds
	if pause < 100000 {
		pause = 0
	}
	val := time.Duration(pause)
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "COUNT", "stat": mem.NumGC})
	ms = append(ms, map[string]interface{}{"value": "TOTAL_TIME", "stat": val.Seconds()})

	for _, val := range mem.PauseNs {
		if val > MAXGC {
			MAXGC = val
		}
	}
	if MAXGC < 100000 {
		MAXGC = 0
	}
	maxtime := time.Duration(MAXGC)

	ms = append(ms, map[string]interface{}{"value": "MAX", "stat": maxtime.Seconds()})

	return createMeaturement(ms)
}

func memoryComitted() []Measurement {
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": mem.HeapSys})
	return createMeaturement(ms)
}

func memoryMax() []Measurement {
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": mem.StackSys})
	return createMeaturement(ms)
}

func cpuCount() []Measurement {
	cpuCount := runtime.NumCPU()

	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": cpuCount})
	return createMeaturement(ms)
}

func buffer() []Measurement {
	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": 0})
	return createMeaturement(ms)
}

func cpuUsage() []Measurement {
	usage, _ := cpu.Percent(0, false)

	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": usage[0]})
	return createMeaturement(ms)
}

func processsCPU() []Measurement {
	p := process.Process{Pid:int32(os.Getpid())}
	usage, _ := p.CPUPercent()

	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": usage})
	return createMeaturement(ms)
}

func processsThreads() []Measurement {
	//usage, _ := p.Threads()

	ms := makeMs()
	ms = append(ms, map[string]interface{}{"value": "VALUE", "stat": 0})
	return createMeaturement(ms)
}


// ###################################################### //
func runMetric(f func() []Measurement) []Measurement {
	return f()
}

func makeMs() []map[string]interface{} {
	return make([]map[string]interface{}, 0)
}

func createMetric(metric string, values map[string]string) Metric {
	return Metric{Name: metric,
		Description: values["description"],
		BaseUnit:    values["baseUnit"],
	}
}

func createMeaturement(ms []map[string]interface{}) []Measurement {
	measurement := make([]Measurement, 0)
	for _, val := range ms {
		measurement = append(measurement, Measurement{Statistic: val["value"].(string), Value: val["stat"]})
	}
	return measurement
}

var metricValues = map[string]map[string]string{
	"process.uptime":            {"description": "Uptime of the application", "baseUnit": "seconds"},
	"jvm.threads.states":        {"description": "The current number of threads having TERMINATED state", "baseUnit": "threads"},
	"jvm.memory.max":            {"description": "The maximum amount of memory in bytes that can be used for memory management", "baseUnit": "bytes"},
	"jvm.gc.memory.promoted":    {"description": "The maximum amount of memory in bytes that can be used for memory management NOT USED BY GO WILL ALWAYS BE 0", "baseUnit": "bytes"},
	"jvm.memory.used":           {"description": "The amount of memory used", "baseUnit": "bytes"},
	"jvm.gc.max.data.size":      {"description": "Max size of old generation memory pool (NOT USE BY GO) ", "baseUnit": "bytes"},
	"jvm.gc.pause":              {"description": "Time spent in GC pause", "baseUnit": "seconds"},
	"jvm.memory.committed":      {"description": "The amount of memory in bytes that is committed for the Java virtual machine to use", "baseUnit": "bytes"},
	"system.cpu.count":          {"description": "The number of processors available to the Java virtual machine", "baseUnit": ""},
	"jvm.buffer.memory.used":    {"description": "An estimate of the memory that the Java virtual machine is using for this buffer pool", "baseUnit": "bytes"},
	"jvm.threads.daemon":        {"description": "NOT USED BY GO WILL ALWAYS BE 0", "baseUnit": "bytes"},
	"system.cpu.usage":          {"description": "The \"recent cpu usage\" for the whole system", "baseUnit": "bytes"},
	"jvm.gc.memory.allocated":   {"description": "Incremented for an increase in the size of the young generation memory pool after one GC to before the next", "baseUnit": "bytes"},
	"jvm.threads.live":          {"description": "NOT USED BY GO WILL ALWAYS BE 0", "baseUnit": "threads"},
	"jvm.threads.peak":          {"description": "NOT USED BY GO WILL ALWAYS BE 0", "baseUnit": "threads"},
	"process.cpu.usage":         {"description": "The \"recent cpu usage\" for the Java Virtual Machine process", "baseUnit": ""},
	"jvm.gc.live.data.size":     {"description": "Size of old generation memory pool after a full GC", "baseUnit": "bytes"},
	"jvm.buffer.count":          {"description": "An estimate of the number of buffers in the pool", "baseUnit": "buffers"},
	"jvm.buffer.total.capacity": {"description": "An estimate of the total capacity of the buffers in this pool", "baseUnit": "bytes"},
	"process.start.time":        {"description": "Start time of the process since unix epoch.", "baseUnit": "seconds"},
}
