package sysinfo

import (
	"fmt"

	"peekfetch/internal/types"

	"github.com/shirou/gopsutil/v3/mem"
)

// GetMemoryInfo collects detailed memory information
func GetMemoryInfo() types.Section {
	info := make(map[string]string)
	order := []string{}

	v, err := mem.VirtualMemory()
	if err == nil {
		// Total
		info["Total RAM"] = formatBytes(v.Total)
		order = append(order, "Total RAM")

		// Used
		info["Used"] = formatBytes(v.Used)
		order = append(order, "Used")

		// Available
		info["Available"] = formatBytes(v.Available)
		order = append(order, "Available")

		// Free
		info["Free"] = formatBytes(v.Free)
		order = append(order, "Free")

		// Cached
		if v.Cached > 0 {
			info["Cached"] = formatBytes(v.Cached)
			order = append(order, "Cached")
		}

		// Buffers
		if v.Buffers > 0 {
			info["Buffers"] = formatBytes(v.Buffers)
			order = append(order, "Buffers")
		}

		// Shared
		if v.Shared > 0 {
			info["Shared"] = formatBytes(v.Shared)
			order = append(order, "Shared")
		}

		// Usage percentage
		info["Usage"] = fmt.Sprintf("%.1f%%", v.UsedPercent)
		order = append(order, "Usage")
	}

	// Swap information
	s, err := mem.SwapMemory()
	if err == nil && s.Total > 0 {
		info["Swap Total"] = formatBytes(s.Total)
		order = append(order, "Swap Total")

		info["Swap Used"] = formatBytes(s.Used)
		order = append(order, "Swap Used")

		info["Swap Free"] = formatBytes(s.Free)
		order = append(order, "Swap Free")

		info["Swap Usage"] = fmt.Sprintf("%.1f%%", s.UsedPercent)
		order = append(order, "Swap Usage")
	}

	return types.Section{
		Name:     "Memory",
		Expanded: false,
		Data:     info,
		LiveData: true,
		Order:    order,
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
