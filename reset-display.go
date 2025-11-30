package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"divoom-monitor/pixoo"
)

func main() {
	host := flag.String("host", "", "Pixoo 64 device IP address (required)")
	channel := flag.Int("channel", 0, "Channel to switch to: 0=Clock, 1=Cloud, 2=Visualizer, 3=Custom")
	flag.Parse()

	if *host == "" {
		fmt.Println("Error: -host flag is required")
		fmt.Println("\nUsage: ./reset-display -host 192.168.1.140")
		fmt.Println("\nChannels:")
		fmt.Println("  0 = Clock/Faces (default)")
		fmt.Println("  1 = Cloud Channel")
		fmt.Println("  2 = Visualizer")
		fmt.Println("  3 = Custom")
		fmt.Println("\nExample: ./reset-display -host 192.168.1.140 -channel 1")
		flag.Usage()
		os.Exit(1)
	}

	client := pixoo.NewClient(*host)

	log.Println("Resetting display...")

	// Clear any custom drawings
	if err := client.ClearScreen(); err != nil {
		log.Printf("Warning: failed to clear screen: %v", err)
	}

	// Switch to specified channel
	channelNames := map[int]string{
		0: "Clock/Faces",
		1: "Cloud Channel",
		2: "Visualizer",
		3: "Custom",
	}

	log.Printf("Switching to channel %d (%s)...", *channel, channelNames[*channel])
	if err := client.SetChannel(*channel); err != nil {
		log.Fatalf("Failed to set channel: %v", err)
	}

	log.Println("Display reset successfully!")
}
