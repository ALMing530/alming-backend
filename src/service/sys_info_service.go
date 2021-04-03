package service

import (
	"alming_backend/src/engine"
	"alming_backend/src/entity"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var tickers = make(map[string]*time.Ticker)

func GetSysInfo() entity.SysInfo {
	cpu, _ := cpu.Percent(time.Second, false)
	memInfo, _ := mem.VirtualMemory()
	return entity.SysInfo{
		Cpu:  cpu[0],
		Mem:  memInfo.UsedPercent,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
}
func PushUsageData(ws *engine.Ws) {
	ws.Conn.WriteJSON(GetSysInfo())
	ticker := time.NewTicker(time.Second * 5)
	tickers[ws.ConnID] = ticker
	go func() {
		for range ticker.C {
			// fmt.Println(s.Format("2006-01-02 15:04:05"))
			// ws.WriteBroadcast([]byte("Broadcast"))
			// ws.WriteByID(ws.ConnID, []byte("Point to Point"))
			// ws.Conn.WriteMessage(websocket.TextMessage, []byte("test data"))
			ws.Conn.WriteJSON(GetSysInfo())
		}
	}()
}
func StopPush(ws *engine.Ws) {
	tickers[ws.ConnID].Stop()
}
