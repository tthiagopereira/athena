package monitoring

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"monitor-system/logger"
	"time"
)

type SystemMonitor struct {
	Interval time.Duration
}

func NewSystemMonitor(interval time.Duration) *SystemMonitor {
	return &SystemMonitor{
		Interval: interval,
	}
}

func (monitor *SystemMonitor) Start(log *logger.Logger) {
	for {
		memInfo, _ := mem.VirtualMemory()
		diskInfo, _ := disk.Usage("/")
		loadInfo, _ := load.Avg()

		info := map[string]string{
			"espaco_total_disco_gb":          fmt.Sprintf("%.2f", float64(diskInfo.Total)/1024/1024/1024),
			"espaco_usado_disco_gb":          fmt.Sprintf("%.2f", float64(diskInfo.Used)/1024/1024/1024),
			"memoria_total_gb":               fmt.Sprintf("%.2f", float64(memInfo.Total)/1024/1024/1024),
			"memoria_usada_gb":               fmt.Sprintf("%.2f", float64(memInfo.Used)/1024/1024/1024),
			"porcentagem_espaco_usado_disco": fmt.Sprintf("%.2f", diskInfo.UsedPercent),
			"porcentagem_memoria_usada":      fmt.Sprintf("%.2f", memInfo.UsedPercent),
			"utilizacao_cpu_1_min":           fmt.Sprintf("%.2f", loadInfo.Load1),
		}

		jsonData, err := json.Marshal(info)
		if err != nil {
			fmt.Println("Erro ao converter para JSON:", err)
			return
		}

		fmt.Println(string(jsonData))

		if err := log.SendLog(jsonData); err != nil {
			fmt.Println("Erro ao enviar JSON para o Elasticsearch:", err)
		}

		time.Sleep(monitor.Interval)
	}
}
