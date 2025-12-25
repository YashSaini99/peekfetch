# PeekFetch

⚡ An interactive system information TUI built with Go and Bubble Tea

## Features

- 🎯 **Keyboard Navigation** - Navigate through sections with arrow keys
- 📊 **Live Updates** - Real-time CPU and memory monitoring (500ms refresh)
- 🎨 **Beautiful UI** - Styled with Lipgloss for a modern terminal experience
- 🚀 **Fast & Lightweight** - Single binary with no dependencies
- 📦 **Multiple Sections** - System, CPU, Memory, Disk, and Network information

## Installation

### Arch Linux (AUR)

```bash
# Using yay
yay -S peekfetch

# Using paru
paru -S peekfetch

# Manual installation
git clone https://aur.archlinux.org/peekfetch.git
cd peekfetch
makepkg -si
```

### Fedora (COPR)

```bash
# Enable the repository
sudo dnf copr enable YashSaini99/peekfetch

# Install
sudo dnf install peekfetch
```

### Build from Source

#### Prerequisites

- Go 1.21 or higher
- Linux (x86_64)

#### Steps

```bash
# Clone the repository
git clone https://github.com/YashSaini99/peekfetch.git
cd peekfetch

# Install dependencies
go mod tidy

# Build the binary
go build -o peekfetch ./cmd/peekfetch

# Run
./peekfetch

# Optional: Install system-wide
sudo install -Dm755 peekfetch /usr/local/bin/peekfetch
```

## Usage

Simply run the binary:

```bash
./peekfetch
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `↑` / `k` | Move selection up |
| `↓` / `j` | Move selection down |
| `Enter` / `Space` | Expand/collapse selected section |
| `L` | Toggle live mode (updates CPU & Memory) |
| `Q` / `Ctrl+C` | Quit application |

## Sections

### System
- Hostname
- User
- Operating System
- Kernel Version
- Architecture
- Uptime
- Boot Time
- Shell (with version)
- Terminal
- Desktop Environment / Window Manager
- Display Server
- Load Average (1m, 5m, 15m)
- Process Count

### CPU
- Model Name
- Vendor ID
- CPU Family & Model ID
- Stepping
- Frequency (GHz)
- Cache Size
- CPU Features/Flags
- Physical Cores
- Logical Cores
- Threads per Core
- Temperature (if available)
- Usage (live updates in live mode)

### Memory
- Total RAM
- Used Memory
- Available Memory
- Free Memory
- Cached Memory
- Buffers
- Shared Memory
- Memory Usage Percentage (live updates in live mode)
- Swap Total
- Swap Used
- Swap Free
- Swap Usage Percentage

### Disk
- Multiple Partitions Support
- For each partition:
  - Mount Point
  - Device Name
  - Filesystem Type
  - Total Space
  - Used Space
  - Free Space
  - Usage Percentage
  - Inode Information

### Network
- Multiple Interfaces Support
- For each interface:
  - Interface Name
  - MAC Address
  - Status Flags
  - MTU
  - IPv4 Address & Subnet
  - IPv6 Address (if available)
- Network Statistics:
  - Total Bytes Sent/Received
  - Total Packets Sent/Received
  - Errors and Drops

## Live Mode

Press `L` to enable live mode. When active:
- CPU usage updates every 500ms
- Memory statistics update every 500ms
- A **[LIVE]** badge appears in the header

## Technical Details

- **Language**: Go
- **TUI Framework**: Bubble Tea
- **Styling**: Lipgloss
- **System Info**: gopsutil v3
- **Target**: Linux x86_64

## Project Structure

```
peekfetch/
├── cmd/
│   └── peekfetch/
│       └── main.go         # Application entry point
├── internal/
│   ├── sysinfo/           # System information gathering
│   │   ├── system.go      # System details (OS, kernel, shell, etc.)
│   │   ├── cpu.go         # CPU information (model, cores, temp, etc.)
│   │   ├── memory.go      # Memory and swap information
│   │   ├── disk.go        # Disk partitions and usage
│   │   └── network.go     # Network interfaces and stats
│   ├── ui/                # User interface
│   │   ├── model.go       # Bubble Tea model
│   │   ├── view.go        # View rendering logic
│   │   └── styles.go      # Lipgloss styling
│   └── types/
│       └── section.go     # Section data structure
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── Makefile               # Build automation
└── README.md              # This file
```

## Building for Release

```bash
# Build optimized binary
go build -ldflags="-s -w" -o peekfetch

# Optional: Compress with UPX
upx peekfetch
```

## License

[MIT](https://github.com/YashSaini99/peekfetch/blob/main/LICENSE)

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [gopsutil](https://github.com/shirou/gopsutil) - System information

---

