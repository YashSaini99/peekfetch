package sysinfo

import (
	"fmt"
	"net"
	"strings"

	"peekfetch/internal/types"

	gopsutilnet "github.com/shirou/gopsutil/v3/net"
)

// GetNetworkInfo collects detailed network information
func GetNetworkInfo() types.Section {
	treeData := []types.TreeItem{}

	interfaces, err := gopsutilnet.Interfaces()
	if err != nil {
		return types.Section{
			Name:     "Network",
			Expanded: false,
			TreeData: treeData,
			LiveData: false,
			UseTree:  true,
		}
	}

	interfaceNum := 1
	for _, iface := range interfaces {
		// Skip loopback
		if strings.HasPrefix(iface.Name, "lo") {
			continue
		}

		// Get Go net interface for addresses
		goIface, err := net.InterfaceByName(iface.Name)
		if err != nil {
			continue
		}

		addrs, err := goIface.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}

		item := types.TreeItem{
			Name:     fmt.Sprintf("Interface %d", interfaceNum),
			Children: make(map[string]string),
			Order:    []string{},
		}

		item.Children["Name"] = iface.Name
		item.Order = append(item.Order, "Name")

		// MAC Address
		if len(iface.HardwareAddr) > 0 {
			item.Children["MAC"] = iface.HardwareAddr
			item.Order = append(item.Order, "MAC")
		}

		// Flags/Status
		if len(iface.Flags) > 0 {
			item.Children["Status"] = strings.Join(iface.Flags, ", ")
			item.Order = append(item.Order, "Status")
		}

		// MTU
		if iface.MTU > 0 {
			item.Children["MTU"] = fmt.Sprintf("%d", iface.MTU)
			item.Order = append(item.Order, "MTU")
		}

		// IP addresses
		ipv4Found := false
		ipv6Found := false

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipNet.IP.To4() != nil && !ipv4Found {
				item.Children["IPv4"] = ipNet.IP.String()
				item.Order = append(item.Order, "IPv4")

				// Subnet mask
				maskSize, _ := ipNet.Mask.Size()
				item.Children["Subnet"] = fmt.Sprintf("/%d", maskSize)
				item.Order = append(item.Order, "Subnet")

				ipv4Found = true
			} else if ipNet.IP.To4() == nil && !ipNet.IP.IsLoopback() && !ipv6Found {
				item.Children["IPv6"] = ipNet.IP.String()
				item.Order = append(item.Order, "IPv6")
				ipv6Found = true
			}
		}

		treeData = append(treeData, item)
		interfaceNum++
	}

	// Get network statistics as a separate tree item
	if stats := getNetworkStats(); len(stats) > 0 {
		statsItem := types.TreeItem{
			Name:     "Statistics",
			Children: stats,
			Order:    []string{"Total Bytes Sent", "Total Bytes Recv", "Total Packets Sent", "Total Packets Recv"},
		}
		if _, ok := stats["Errors"]; ok {
			statsItem.Order = append(statsItem.Order, "Errors")
		}
		if _, ok := stats["Drops"]; ok {
			statsItem.Order = append(statsItem.Order, "Drops")
		}
		treeData = append(treeData, statsItem)
	}

	return types.Section{
		Name:     "Network",
		Expanded: false,
		TreeData: treeData,
		LiveData: false,
		UseTree:  true,
	}
}

func getNetworkStats() map[string]string {
	stats := make(map[string]string)

	ioCounters, err := gopsutilnet.IOCounters(false)
	if err != nil || len(ioCounters) == 0 {
		return stats
	}

	total := ioCounters[0]

	stats["Total Bytes Sent"] = formatBytes(total.BytesSent)
	stats["Total Bytes Recv"] = formatBytes(total.BytesRecv)
	stats["Total Packets Sent"] = fmt.Sprintf("%d", total.PacketsSent)
	stats["Total Packets Recv"] = fmt.Sprintf("%d", total.PacketsRecv)

	if total.Errin > 0 || total.Errout > 0 {
		stats["Errors"] = fmt.Sprintf("In: %d, Out: %d", total.Errin, total.Errout)
	}

	if total.Dropin > 0 || total.Dropout > 0 {
		stats["Drops"] = fmt.Sprintf("In: %d, Out: %d", total.Dropin, total.Dropout)
	}

	return stats
}
