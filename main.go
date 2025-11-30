package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"divoom-monitor/metrics"
	"divoom-monitor/pixoo"
)

func main() {
	// Parse command line flags
	host := flag.String("host", "", "Pixoo 64 device IP address (required)")
	interval := flag.Int("interval", 5, "Update interval in seconds")
	brightness := flag.Int("brightness", 50, "Screen brightness (0-100)")
	textOnly := flag.Bool("text", false, "Use text-only mode (faster, less detailed)")
	flag.Parse()

	if *host == "" {
		fmt.Println("Error: -host flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create Pixoo client
	client := pixoo.NewClient(*host)

	// Set brightness
	if err := client.SetBrightness(*brightness); err != nil {
		log.Printf("Warning: failed to set brightness: %v", err)
	}

	// Switch to Custom channel (channel 3) so our drawings appear
	log.Println("Switching to Custom channel...")
	if err := client.SetChannel(3); err != nil {
		log.Printf("Warning: failed to set channel: %v", err)
	}

	// Create metrics collector
	collector := metrics.NewCollector()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Printf("Starting Divoom monitor on %s (update every %ds)", *host, *interval)
	log.Println("Press Ctrl+C to exit")

	// Main loop
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	// Initial update
	if err := updateDisplay(client, collector, *textOnly); err != nil {
		log.Printf("Error updating display: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := updateDisplay(client, collector, *textOnly); err != nil {
				log.Printf("Error updating display: %v", err)
			}
		case <-sigChan:
			log.Println("Shutting down...")
			return
		}
	}
}

func updateDisplay(client *pixoo.Client, collector *metrics.Collector, textOnly bool) error {
	// Collect system metrics
	m, err := collector.Collect()
	if err != nil {
		return fmt.Errorf("collect metrics: %w", err)
	}

	// Log metrics
	log.Println(m.String())

	// Text-only mode (faster, uses Pixoo's built-in text rendering)
	if textOnly {
		text := fmt.Sprintf("CPU:%.0f%% MEM:%.0f%% %.1fG",
			m.CPUPercent, m.MemoryPercent, m.MemoryUsedGB)
		if err := client.DrawText(text, 255, 255, 255); err != nil {
			return fmt.Errorf("draw text: %w", err)
		}
		return nil
	}

	// Image mode (slower but with graphics)
	img := pixoo.CreateImage()

	// Colors
	bgColor := color.RGBA{0, 0, 0, 255}        // Black background
	cpuColor := color.RGBA{0, 255, 0, 255}     // Green for CPU
	memColor := color.RGBA{0, 150, 255, 255}   // Blue for Memory
	textColor := color.RGBA{255, 255, 255, 255} // White text

	// Fill background
	pixoo.FillRect(img, 0, 0, 64, 64, bgColor)

	// Draw title
	pixoo.DrawTextOnImage(img, "CPU:", 2, 2, textColor)
	cpuText := fmt.Sprintf("%2.0f%%", m.CPUPercent)
	pixoo.DrawTextOnImage(img, cpuText, 30, 2, textColor)

	// Draw CPU bar
	cpuBarWidth := int((m.CPUPercent / 100.0) * 60)
	pixoo.FillRect(img, 2, 12, 2+cpuBarWidth, 16, cpuColor)

	// Draw memory info
	pixoo.DrawTextOnImage(img, "MEM:", 2, 20, textColor)
	memText := fmt.Sprintf("%2.0f%%", m.MemoryPercent)
	pixoo.DrawTextOnImage(img, memText, 30, 20, textColor)

	// Draw memory bar
	memBarWidth := int((m.MemoryPercent / 100.0) * 60)
	pixoo.FillRect(img, 2, 30, 2+memBarWidth, 34, memColor)

	// Draw memory usage in GB
	memGBText := fmt.Sprintf("%.1fG", m.MemoryUsedGB)
	pixoo.DrawTextOnImage(img, memGBText, 2, 38, textColor)

	// Draw network stats (if available)
	if m.NetRecvMB > 0 || m.NetSentMB > 0 {
		netText := fmt.Sprintf("%.1fM", m.NetRecvMB)
		pixoo.DrawTextOnImage(img, netText, 2, 50, textColor)
	}

	// Send image to display
	log.Println("Sending image to display...")
	if err := client.DrawImage(img); err != nil {
		return fmt.Errorf("draw image: %w", err)
	}
	log.Println("Image sent successfully")

	return nil
}
