package sysinfo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"peekfetch/internal/types"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
)

// GetSystemInfo collects detailed system information
func GetSystemInfo() types.Section {
	info := make(map[string]string)
	order := []string{}

	// Hostname
	hostname, _ := os.Hostname()
	info["Hostname"] = hostname
	order = append(order, "Hostname")

	// User
	if user := os.Getenv("USER"); user != "" {
		info["User"] = user
		order = append(order, "User")
	}

	// OS
	hostInfo, _ := host.Info()
	info["OS"] = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
	order = append(order, "OS")

	// Kernel
	info["Kernel"] = hostInfo.KernelVersion
	order = append(order, "Kernel")

	// Architecture
	info["Architecture"] = runtime.GOARCH
	order = append(order, "Architecture")

	// Uptime
	uptime := time.Duration(hostInfo.Uptime) * time.Second
	info["Uptime"] = formatDuration(uptime)
	order = append(order, "Uptime")

	// Boot time
	bootTime := time.Unix(int64(hostInfo.BootTime), 0)
	info["Boot Time"] = bootTime.Format("2006-01-02 15:04:05")
	order = append(order, "Boot Time")

	// Shell
	if shell := os.Getenv("SHELL"); shell != "" {
		// Extract just the shell name
		parts := strings.Split(shell, "/")
		shellName := parts[len(parts)-1]

		// Try to get version
		if version := getShellVersion(shellName); version != "" {
			info["Shell"] = fmt.Sprintf("%s %s", shellName, version)
		} else {
			info["Shell"] = shellName
		}
		order = append(order, "Shell")
	}

	// Terminal
	if term := os.Getenv("TERM"); term != "" {
		info["Terminal"] = term
		order = append(order, "Terminal")
	}

	// Desktop Environment / Window Manager
	if de := getDesktopEnvironment(); de != "" {
		info["Desktop"] = de
		order = append(order, "Desktop")
	}

	// Display info
	if display := os.Getenv("DISPLAY"); display != "" {
		info["Display"] = display
		order = append(order, "Display")
	}

	// Load average
	if avg, err := load.Avg(); err == nil {
		info["Load Average"] = fmt.Sprintf("%.2f, %.2f, %.2f", avg.Load1, avg.Load5, avg.Load15)
		order = append(order, "Load Average")
	}

	// Number of processes
	if procs := countProcesses(); procs > 0 {
		info["Processes"] = fmt.Sprintf("%d", procs)
		order = append(order, "Processes")
	}

	return types.Section{
		Name:     "System",
		Expanded: false,
		Data:     info,
		LiveData: false,
		Order:    order,
	}
}

func getShellVersion(shell string) string {
	var cmd *exec.Cmd

	switch shell {
	case "bash":
		cmd = exec.Command("bash", "--version")
	case "zsh":
		cmd = exec.Command("zsh", "--version")
	case "fish":
		cmd = exec.Command("fish", "--version")
	default:
		return ""
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		// Extract version from first line
		parts := strings.Fields(lines[0])
		for i, part := range parts {
			if strings.Contains(part, ".") && i > 0 {
				return strings.TrimSpace(part)
			}
		}
	}

	return ""
}

func getDesktopEnvironment() string {
	// Try various environment variables
	envVars := []string{
		"XDG_CURRENT_DESKTOP",
		"DESKTOP_SESSION",
		"XDG_SESSION_DESKTOP",
	}

	for _, envVar := range envVars {
		if val := os.Getenv(envVar); val != "" {
			return val
		}
	}

	// Try to detect window manager
	wms := []string{"i3", "sway", "bspwm", "awesome", "dwm", "xmonad"}
	for _, wm := range wms {
		if _, err := exec.LookPath(wm); err == nil {
			// Check if it's running
			cmd := exec.Command("pgrep", "-x", wm)
			if err := cmd.Run(); err == nil {
				return wm
			}
		}
	}

	return ""
}

func countProcesses() int {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 4 {
			// Format is "load1 load5 load15 running/total pid"
			if strings.Contains(fields[3], "/") {
				parts := strings.Split(fields[3], "/")
				if len(parts) == 2 {
					var total int
					fmt.Sscanf(parts[1], "%d", &total)
					return total
				}
			}
		}
	}

	return 0
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}

	return strings.Join(parts, " ")
}
