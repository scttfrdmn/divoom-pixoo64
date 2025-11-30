#!/bin/bash

# Pixoo 64 Monitor using curl (bypasses macOS firewall issues)

HOST="${1:-}"
INTERVAL="${2:-5}"

if [ -z "$HOST" ]; then
    echo "Usage: $0 <pixoo-ip> [interval-seconds]"
    echo "Example: $0 192.168.1.140 5"
    exit 1
fi

URL="http://$HOST/post"

echo "Starting Pixoo 64 monitor on $HOST (update every ${INTERVAL}s)"
echo "Press Ctrl+C to exit"
echo ""

# Function to get CPU usage
get_cpu() {
    top -l 2 -n 0 -F | grep "CPU usage" | tail -1 | awk '{print $3}' | sed 's/%//'
}

# Function to get memory usage
get_memory() {
    vm_stat | perl -ne '/page size of (\d+)/ and $size=$1; /Pages\s+([^:]+)[^\d]+(\d+)/ and printf("%-16s % 16.2f Mi\n", "$1:", $2 * $size / 1048576);' | grep "active\|wired" | awk '{sum+=$2} END {printf "%.1f", sum/1024}'
}

# Function to send image to Pixoo
send_image() {
    local cpu=$1
    local mem=$2

    # Create a simple visualization with bars
    # For now, we'll just send text
    local text="CPU:${cpu}% MEM:${mem}G"

    local json="{\"Command\":\"Draw/SendHttpText\",\"TextId\":1,\"x\":0,\"y\":24,\"dir\":0,\"font\":2,\"TextWidth\":64,\"speed\":0,\"TextString\":\"$text\",\"color\":\"#FFFFFF\",\"align\":2}"

    curl -s -X POST "$URL" \
         -H "Content-Type: application/json" \
         -d "$json" > /dev/null
}

# Main loop
while true; do
    CPU=$(get_cpu)
    MEM=$(get_memory)

    echo "$(date '+%Y/%m/%d %H:%M:%S') CPU: ${CPU}% | MEM: ${MEM}GB"

    send_image "$CPU" "$MEM"

    sleep "$INTERVAL"
done
