# Troubleshooting Divoom Monitor

## "no route to host" Error on macOS

If you see errors like:
```
dial tcp 192.168.1.140:80: connect: no route to host
```

This is caused by **Tailscale** or other VPN software blocking the connection, NOT the macOS firewall.

### Quick Test

Run these commands to verify Tailscale is the issue:
```bash
# Check if Tailscale is running
ps aux | grep tailscale | grep -v grep

# Check VPN interfaces
ifconfig | grep utun
```

If you see tailscale processes or multiple `utun` interfaces, that's the issue.

### Solutions

#### Option 1: Temporarily Disable Tailscale (Recommended for testing)

```bash
# Stop Tailscale
tailscale down

# Run your monitor
./divoom-monitor -host 192.168.1.140

# When done, restart Tailscale
tailscale up
```

#### Option 2: Use the curl-based script

The curl script bypasses this issue entirely:
```bash
./pixoo-curl.sh 192.168.1.140 5
```

#### Option 3: Run with sudo

Sometimes elevated privileges bypass VPN restrictions:
```bash
sudo ./divoom-monitor -host 192.168.1.140
```

#### Option 4: Enable LAN access in Tailscale

1. Click the Tailscale icon in your menu bar
2. Open "Preferences" or "Settings"
3. Look for one of these options:
   - "Allow LAN access"
   - "Local network access"
   - "Bypass Tailscale for local addresses"
4. Enable it and try again

#### Option 5: Add to Tailscale ACL

If you have admin access to your Tailscale network, you can add ACL rules to allow local network access.

### Verification

To verify your Pixoo is reachable and the issue is only with Go:

```bash
# Should work - ping test
ping 192.168.1.140

# Should work - port test
nc -zv 192.168.1.140 80

# Should work - curl test
curl http://192.168.1.140/post

# Fails - Go binary (if Tailscale is blocking)
./divoom-monitor -host 192.168.1.140
```

## Other Network Issues

### Device Not Found

If ping fails:
```bash
# Find your Pixoo IP address
# Option 1: Check the Divoom app on your phone

# Option 2: Scan your network
nmap -sn 192.168.1.0/24

# Option 3: Check ARP table
arp -a | grep -i divoom
```

### Port 80 Not Responding

If nc test fails but ping works:
- Make sure the Pixoo is not in sleep mode
- Restart the Pixoo device
- Check if another app is using the device
- Try accessing from the Divoom mobile app first

### Slow Updates

If the display updates but slowly:
- Reduce the update interval: `-interval 10`
- Check your WiFi signal strength
- Restart your router
- Move the Pixoo closer to your WiFi access point

## Performance Issues

### High CPU Usage

The application should use minimal CPU. If you see high usage:
- Increase the interval: `-interval 10` or higher
- Check if other processes are consuming resources
- Monitor with: `top -pid $(pgrep divoom-monitor)`

### Memory Leaks

If memory usage grows over time:
- Restart the application periodically
- Report the issue with logs: `./divoom-monitor -host IP 2>&1 | tee debug.log`

## Display Issues

### Display Shows Garbled Image

- The Pixoo might be in the wrong mode
- Open the Divoom app and switch to "Custom" or "Cloud Channel" mode
- Restart the application

### Text Not Readable

The included font is 5x7 pixels and limited:
- Supported characters: 0-9, basic letters, % : space
- Modify `pixoo/client.go` `getCharPattern()` to add more characters
- Or use the Pixoo's built-in text rendering with `DrawText()`

### Colors Look Wrong

- Check the Pixoo's brightness setting
- Adjust brightness: `-brightness 80`
- Modify colors in `main.go` lines 112-114

## Development Issues

### Build Errors

```bash
# Clean and rebuild
go clean
go mod tidy
go build -o divoom-monitor
```

### Import Errors

```bash
# Download dependencies
go mod download
go get github.com/shirou/gopsutil/v3
```

### Cannot Find Device

Create a simple test script:
```bash
# test-device.sh
#!/bin/bash
HOST=$1
echo "Testing $HOST..."
curl -X POST http://$HOST/post \
  -H "Content-Type: application/json" \
  -d '{"Command":"Channel/GetIndex"}'
```

```bash
chmod +x test-device.sh
./test-device.sh 192.168.1.140
```

If this works, the device is fine and the issue is with your setup.

## Still Having Issues?

1. Check the main README.md for basic setup instructions
2. Verify your Pixoo works with the official Divoom app
3. Try the curl script as a workaround
4. Create an issue with:
   - Your macOS version: `sw_vers`
   - Tailscale status: `tailscale status | head -5`
   - Error messages from the application
   - Output of the verification tests above

## Working Configurations

### Confirmed Working

- macOS 13+ with Tailscale disabled
- macOS with no VPN software
- Linux (all distributions)
- Using the curl-based script with Tailscale enabled

### Known Issues

- Tailscale blocks Go binaries from local network access
- Some corporate VPNs block local network entirely
- Some firewall software requires manual approval
- Docker Desktop networking can interfere
