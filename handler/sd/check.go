package sd

import (
	"encoding/json"
	"fmt"
	"net/http"

	v "gopa/pkg/version"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// HealthCheck check service status
// @Summary      HealthCheck shows `OK` as the ping-pong result.
// @Description  HealthCheck
// @Tags         sd
// @Accept       json
// @Success      200  {object}  handler.Response
// @Router       /sd/health [get]
func HealthCheck(c *gin.Context) {
	message := "OK"
	c.String(http.StatusOK, "\n"+message)

}

//VersionCheck show the version info as Running status
// @Summary      VersionCheck show the version info as Running status
// @Description  versionCheck
// @Tags         sd
// @Accept       json
// @Success      200  {object}  handler.Response
// @Router       /sd/version [get]
func VersionCheck(c *gin.Context) {
	v := v.Get()
	fmt.Println("version:", v)
	marshalled, _ := json.MarshalIndent(&v, "", "  ")
	fmt.Println("message:", string(marshalled))
	c.String(http.StatusOK, "\n"+string(marshalled))
}

//DiskCheck checks the disk usage.
// @Summary      DiskCheck checks the disk usage.
// @Description  DiskCheck
// @Tags         sd
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.Response
// @Router       /sd/disk [get]
func DiskCheck(c *gin.Context) {
	u, _ := disk.Usage("/")

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusOK
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.String(status, "\n"+message)
}

//CPUCheck checks the cpu usage.
// @Summary      CPUCheck checks the cpu usage.
// @Description  CPUCheck
// @Tags         sd
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.Response
// @Router       /sd/cpu [get]
func CPUCheck(c *gin.Context) {
	cores, _ := cpu.Counts(false)

	a, _ := load.Avg()
	l1 := a.Load1
	l5 := a.Load5
	l15 := a.Load15

	status := http.StatusOK
	text := "OK"

	if l5 >= float64(cores-1) {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if l5 >= float64(cores-2) {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Cores: %d", text, l1, l5, l15, cores)
	c.String(status, "\n"+message)
}

//RAMCheck checks the disk usage.
// @Summary      RAMCheck checks the disk usage.
// @Description  RAMCheck
// @Tags         sd
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.Response
// @Router       /sd/ram [get]
func RAMCheck(c *gin.Context) {
	u, _ := mem.VirtualMemory()

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.String(status, "\n"+message)
}
