package service

import (
	"alming_backend/src/entity"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetSysInfo() entity.SysInfo {
	cpu, _ := cpu.Percent(time.Second, false)
	memInfo, _ := mem.VirtualMemory()
	return entity.SysInfo{
		Cpu: cpu[0],
		Mem: memInfo.UsedPercent,
	}
}
