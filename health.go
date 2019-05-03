package fairway

import (
	"github.com/shirou/gopsutil/disk"
)

type healthData struct {
	Status  string                 `json:"status"`
	Details map[string]interface{} `json:"details"`
}

type diskSpace struct {
	Status  string            `json:"status"`
	Details map[string]uint64 `json:"details"`
}

func health() ([]byte, error) {

	health := healthData{Status: "UP"}

	det := map[string]interface{}{}

	ds, err := getDiskSpace()
	if err != nil {
		return []byte(""), err
	}
	det["diskSpace"] = ds

	health.Details = det

	return toJSON(health), nil
}

func getDiskSpace() (diskSpace, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return diskSpace{}, err
	}
	diskSpace := diskSpace{}
	d := map[string]uint64{}

	d["free"] = diskStat.Free
	d["used"] = diskStat.Used
	d["total"] = diskStat.Total
	d["threashold"] = 10485760

	diskSpace.Details = d
	return diskSpace, nil
}
