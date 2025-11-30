# Divoom Pixoo 64 System Monitor

A Go application that displays real-time system metrics on your Divoom Pixoo 64 pixel display.

## Features

- **CPU Usage**: Visual bar graph showing current CPU utilization
- **Memory Usage**: Visual bar graph and GB usage display
- **Network Stats**: Download speeds in MB/s
- **Auto-refresh**: Configurable update interval
- **Customizable**: Brightness control and update frequency

## Prerequisites

- Go 1.21 or higher
- Divoom Pixoo 64 connected to your local network
- Device IP address

## Finding Your Pixoo 64 IP Address

1. Open the Divoom app on your phone
2. Go to device settings
3. Look for "Device IP" or check your router's DHCP client list
4. Alternatively, scan your network:
   ```bash
   # On macOS/Linux
   arp -a | grep divoom

   # Or use nmap
   nmap -sn 192.168.1.0/24
   ```

## Installation

1. Clone or download this repository

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o divoom-monitor
   ```

## ⚠️ Important: Tailscale/VPN Users

**If you have Tailscale or other VPN software running**, the Go binary will be blocked from accessing local network devices by default.

**Quick solutions:**
```bash
# Option 1: Temporarily disable Tailscale
tailscale down
./divoom-monitor -host 192.168.1.140
tailscale up

# Option 2: Use the curl-based script (works with Tailscale)
./pixoo-curl.sh 192.168.1.140 5

# Option 3: Try with sudo
sudo ./divoom-monitor -host 192.168.1.140
```

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for detailed solutions.

## Usage

Basic usage with required IP address:
```bash
./divoom-monitor -host 192.168.1.100
```

With all options:
```bash
./divoom-monitor -host 192.168.1.100 -interval 5 -brightness 70
```

### Command Line Options

- `-host` (required): IP address of your Pixoo 64 device
- `-interval`: Update interval in seconds (default: 5)
- `-brightness`: Screen brightness 0-100 (default: 50)

### Example

```bash
# Update every 3 seconds with 80% brightness
./divoom-monitor -host 192.168.1.150 -interval 3 -brightness 80
```

## Display Layout

The 64x64 pixel display shows:

```
┌────────────────┐
│ CPU:  XX%      │
│ ████████░░░░░  │ (Green bar)
│                │
│ MEM:  XX%      │
│ ████████░░░░░  │ (Blue bar)
│ X.XG           │
│                │
│ ↓ X.XM         │ (Network down)
└────────────────┘
```

## Project Structure

```
divoom-monitor/
├── main.go              # Main application
├── pixoo/
│   └── client.go        # Pixoo 64 API client
├── metrics/
│   └── collector.go     # System metrics collector
├── go.mod
└── README.md
```

## API Documentation

### Pixoo Client

The `pixoo` package provides methods to control the Pixoo 64:

```go
client := pixoo.NewClient("192.168.1.100")

// Set brightness (0-100)
client.SetBrightness(70)

// Clear screen
client.ClearScreen()

// Draw custom image (64x64)
img := pixoo.CreateImage()
client.DrawImage(img)

// Draw text
client.DrawText("Hello World", 255, 255, 255)
```

### Metrics Collector

The `metrics` package collects system information:

```go
collector := metrics.NewCollector()
m, err := collector.Collect()

// Access metrics
fmt.Printf("CPU: %.1f%%\n", m.CPUPercent)
fmt.Printf("Memory: %.1f%%\n", m.MemoryPercent)
fmt.Printf("Network: ↓%.2f MB/s\n", m.NetRecvMB)
```

## Customization

### Changing Colors

Edit `main.go` in the `updateDisplay()` function:

```go
cpuColor := color.RGBA{0, 255, 0, 255}     // Green
memColor := color.RGBA{0, 150, 255, 255}   // Blue
textColor := color.RGBA{255, 255, 255, 255} // White
```

### Adding More Metrics

1. Add collection logic in `metrics/collector.go`
2. Update the `SystemMetrics` struct
3. Modify display logic in `main.go`

### Custom Display Layout

The `updateDisplay()` function in `main.go` controls the layout. Modify it to:
- Change text positions
- Add more bars or graphs
- Display different information

## Troubleshooting

### macOS Firewall Blocking Connection (Most Common)

If you see `dial tcp: connect: no route to host` errors, this is macOS's application firewall blocking the unsigned Go binary.

**Solution 1 - Allow through firewall:**
1. Run the application: `./divoom-monitor -host YOUR_IP`
2. macOS will prompt "Do you want the application to accept incoming network connections?"
3. Click "Allow"

**Solution 2 - Use the curl-based script (bypasses firewall):**
```bash
./pixoo-curl.sh 192.168.1.140 5
```

**Solution 3 - Manually allow in firewall:**
1. Open System Settings > Network > Firewall (or Security & Privacy > Firewall)
2. Click "Firewall Options" or "Options"
3. Click the + button and add the `divoom-monitor` binary
4. Set it to "Allow incoming connections"

**Solution 4 - Temporarily disable firewall (not recommended):**
```bash
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate off
# Run your program
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate on
```

### Other Connection Issues

- Verify the Pixoo 64 is on the same network
- Check the IP address is correct: `ping 192.168.1.140`
- Verify port 80 is accessible: `nc -zv 192.168.1.140 80`
- Check if curl works: `curl http://192.168.1.140/post`

### Display Not Updating

- Check the device isn't in another mode (clock, visualization, etc.)
- Restart the application
- Reduce update interval if system is slow

### Permission Errors

On Linux, you may need elevated privileges for system metrics:
```bash
sudo ./divoom-monitor -host 192.168.1.100
```

## Development

### Running from source
```bash
go run main.go -host 192.168.1.100
```

### Running tests
```bash
go test ./...
```

## License

MIT

## Contributing

Feel free to submit issues and pull requests!

## Acknowledgments

- Built with [gopsutil](https://github.com/shirou/gopsutil) for system metrics
- Divoom Pixoo 64 API documentation
