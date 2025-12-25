package sysinfo

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"peekfetch/internal/types"

	"github.com/shirou/gopsutil/v3/cpu"
)

// GetCPUInfo collects detailed CPU information
func GetCPUInfo() types.Section {
	info := make(map[string]string)
	order := []string{}

	cpuInfo, _ := cpu.Info()
	if len(cpuInfo) > 0 {
		// Model
		info["Model"] = cpuInfo[0].ModelName
		order = append(order, "Model")

		// Vendor
		if cpuInfo[0].VendorID != "" {
			info["Vendor"] = cpuInfo[0].VendorID
			order = append(order, "Vendor")
		}

		// CPU family
		if cpuInfo[0].Family != "" {
			info["Family"] = cpuInfo[0].Family
			order = append(order, "Family")
		}

		// CPU model number
		if cpuInfo[0].Model != "" {
			info["Model ID"] = cpuInfo[0].Model
			order = append(order, "Model ID")
		}

		// Stepping
		if cpuInfo[0].Stepping != 0 {
			info["Stepping"] = fmt.Sprintf("%d", cpuInfo[0].Stepping)
			order = append(order, "Stepping")
		}

		// Frequency
		if cpuInfo[0].Mhz > 0 {
			ghz := cpuInfo[0].Mhz / 1000
			info["Frequency"] = fmt.Sprintf("%.2f GHz", ghz)
			order = append(order, "Frequency")
		}

		// Cache size
		if cpuInfo[0].CacheSize > 0 {
			info["Cache Size"] = fmt.Sprintf("%d KB", cpuInfo[0].CacheSize)
			order = append(order, "Cache Size")
		}

		// Flags/Features (show first few)
		if len(cpuInfo[0].Flags) > 0 {
			flags := cpuInfo[0].Flags
			if len(flags) > 10 {
				flags = flags[:10]
			}
			info["Features"] = strings.Join(flags, ", ") + "..."
			order = append(order, "Features")
		}
	}

	// Physical and logical cores
	physicalCores, _ := cpu.Counts(false)
	logicalCores := runtime.NumCPU()

	if physicalCores > 0 {
		info["Physical Cores"] = fmt.Sprintf("%d", physicalCores)
		order = append(order, "Physical Cores")
	}

	info["Logical Cores"] = fmt.Sprintf("%d", logicalCores)
	order = append(order, "Logical Cores")

	// Threads per core
	if physicalCores > 0 && logicalCores > 0 {
		threadsPerCore := logicalCores / physicalCores
		info["Threads/Core"] = fmt.Sprintf("%d", threadsPerCore)
		order = append(order, "Threads/Core")
	}

	// Temperature (if available)
	if temp := getCPUTemperature(); temp != "" {
		info["Temperature"] = temp
		order = append(order, "Temperature")
	}

	// Usage placeholder (will be updated in live mode)
	info["Usage"] = "0.0%"
	order = append(order, "Usage")

	return types.Section{
		Name:     "CPU",
		Expanded: false,
		Data:     info,
		LiveData: true,
		Order:    order,
	}
}

// GetCPUUsage gets current CPU usage (live data)
func GetCPUUsage() string {
	percent, err := cpu.Percent(100*time.Millisecond, false)
	if err != nil || len(percent) == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.1f%%", percent[0])
}

func getCPUTemperature() string {
	// Try different thermal zone files
	thermalPaths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/thermal/thermal_zone1/temp",
		"/sys/class/hwmon/hwmon0/temp1_input",
		"/sys/class/hwmon/hwmon1/temp1_input",
	}

	for _, path := range thermalPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var temp int
		_, err = fmt.Sscanf(string(data), "%d", &temp)
		if err != nil {
			continue
		}

		// Temperature is in millidegrees
		celsius := float64(temp) / 1000.0

		// Only return if reasonable (between 0 and 150°C)
		if celsius > 0 && celsius < 150 {
			return fmt.Sprintf("%.1f°C", celsius)
		}
	}

	return ""
}
