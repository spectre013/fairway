package fairway

import (
	"github.com/shirou/gopsutil/disk"
)


type HealthData struct {
	Status string `json:"status"`
	Details map[string]interface{} `json:"details"`
}

type DiskSpace struct {
	Status string `json:"status"`
	Details map[string]uint64 `json:"details"`
}

func health() ([]byte, error) {

	health := HealthData{Status:"UP"}


	det := map[string]interface{}{}

	ds, err  := getDiskSpace()
	if err != nil {
		return []byte(""), err
	}
	det["diskSpace"] = ds

	health.Details = det

	return toJson(health), nil
}

func getDiskSpace() (DiskSpace,error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return DiskSpace{}, err
	}
	diskSpace := DiskSpace{}
	d := map[string]uint64{}

	d["free"] = diskStat.Free
	d["used"] = diskStat.Used
	d["total"] = diskStat.Total
	d["threashold"] = 10485760

	diskSpace.Details = d
	return diskSpace, nil
}