package logsend

import (
    "time"
    log "github.com/ezotrank/logger"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/load"
)

var (
    SMUpdateInterval = time.Duration(5 * time.Second)
)


func RunSystemMetricsCollect() {
    for {
        metrics := make(map[string]interface{},0)
        metrics["virtual_memory"],_ = mem.VirtualMemory()
        metrics["cpu_info"],_ = load.LoadAvg() 
        log.Infoln(metrics)
        time.Sleep(SMUpdateInterval)
    }
}
