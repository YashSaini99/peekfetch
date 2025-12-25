package sysinfo

import (
	"fmt"
	"os/exec"
	"strings"

	"peekfetch/internal/types"

	"github.com/shirou/gopsutil/v3/disk"
)

// GetDiskInfo collects detailed disk information for all mounted partitions
func GetDiskInfo() types.Section {
	treeData := []types.TreeItem{}

	// Get all partitions
	partitions, err := disk.Partitions(false)
	if err != nil {
		return types.Section{
			Name:     "Disk",
			Expanded: false,
			TreeData: treeData,
			LiveData: false,
			UseTree:  true,
		}
	}

	partNum := 1
	for _, partition := range partitions {
		// Skip special filesystems
		if isSpecialFS(partition.Fstype) {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		item := types.TreeItem{
			Name:     fmt.Sprintf("Partition %d", partNum),
			Children: make(map[string]string),
			Order:    []string{},
		}

		item.Children["Mount"] = partition.Mountpoint
		item.Order = append(item.Order, "Mount")

		item.Children["Device"] = partition.Device
		item.Order = append(item.Order, "Device")

		item.Children["FS Type"] = partition.Fstype
		item.Order = append(item.Order, "FS Type")

		item.Children["Total"] = formatBytes(usage.Total)
		item.Order = append(item.Order, "Total")

		item.Children["Used"] = formatBytes(usage.Used)
		item.Order = append(item.Order, "Used")

		item.Children["Free"] = formatBytes(usage.Free)
		item.Order = append(item.Order, "Free")

		item.Children["Usage"] = fmt.Sprintf("%.1f%%", usage.UsedPercent)
		item.Order = append(item.Order, "Usage")

		// Inodes if available
		if usage.InodesTotal > 0 {
			item.Children["Inodes"] = fmt.Sprintf("%d / %d", usage.InodesUsed, usage.InodesTotal)
			item.Order = append(item.Order, "Inodes")
		}

		treeData = append(treeData, item)
		partNum++
	}

	// Add total disk I/O stats if available
	if ioStats := getDiskIOStats(); ioStats != "" {
		statsItem := types.TreeItem{
			Name:     "I/O Statistics",
			Children: map[string]string{"Status": ioStats},
			Order:    []string{"Status"},
		}
		treeData = append(treeData, statsItem)
	}

	return types.Section{
		Name:     "Disk",
		Expanded: false,
		TreeData: treeData,
		LiveData: false,
		UseTree:  true,
	}
}

func isSpecialFS(fstype string) bool {
	specialFS := []string{
		"tmpfs", "devtmpfs", "devpts", "sysfs", "proc",
		"cgroup", "cgroup2", "pstore", "bpf", "tracefs",
		"debugfs", "hugetlbfs", "mqueue", "configfs",
		"securityfs", "fusectl", "fuse.portal",
	}

	for _, special := range specialFS {
		if fstype == special {
			return true
		}
	}

	return false
}

func getDiskIOStats() string {
	// Try to get basic I/O stats using iostat if available
	cmd := exec.Command("iostat", "-d", "-x", "1", "1")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	// Just return a summary indicator
	if len(lines) > 0 {
		return "Available (use iostat for details)"
	}

	return ""
}
