package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"divoom-monitor/gameoflife"
	"divoom-monitor/pixoo"
)

func main() {
	// Parse command line flags
	host := flag.String("host", "", "Pixoo 64 device IP address (required)")
	pattern := flag.String("pattern", "random", "Starting pattern: random, random-sparse, random-dense, gliders, gosper-gun, pulsar")
	speed := flag.Int("speed", 200, "Update speed in milliseconds")
	brightness := flag.Int("brightness", 70, "Screen brightness (0-100)")
	colorMode := flag.String("color", "age", "Color mode: age, rainbow, fire, ocean, matrix")
	flag.Parse()

	if *host == "" {
		fmt.Println("Error: -host flag is required")
		fmt.Println("\nExample: go run game-of-life.go -host 192.168.1.140")
		fmt.Println("\nPatterns:")
		fmt.Println("  random        - Random starting state (30% density)")
		fmt.Println("  random-sparse - Random sparse (15% density)")
		fmt.Println("  random-dense  - Random dense (50% density)")
		fmt.Println("  gliders       - Multiple glider patterns")
		fmt.Println("  gosper-gun    - Gosper glider gun")
		fmt.Println("  pulsar        - Pulsar oscillator")
		fmt.Println("\nColor modes:")
		fmt.Println("  age     - Color based on cell age")
		fmt.Println("  rainbow - Rainbow gradient")
		fmt.Println("  fire    - Fire colors (red/orange/yellow)")
		fmt.Println("  ocean   - Ocean colors (blue/cyan)")
		fmt.Println("  matrix  - Matrix green")
		flag.Usage()
		os.Exit(1)
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create Pixoo client
	client := pixoo.NewClient(*host)

	// Set brightness
	if err := client.SetBrightness(*brightness); err != nil {
		log.Printf("Warning: failed to set brightness: %v", err)
	}

	// Switch to Custom channel
	log.Println("Switching to Custom channel...")
	if err := client.SetChannel(3); err != nil {
		log.Printf("Warning: failed to set channel: %v", err)
	}

	// Create game
	game := gameoflife.NewGame()
	game.LoadPattern(*pattern)

	log.Printf("Starting Game of Life on %s", *host)
	log.Printf("Pattern: %s, Speed: %dms, Color: %s", *pattern, *speed, *colorMode)
	log.Println("Press Ctrl+C to exit")

	// Track cell ages for coloring
	cellAges := make([][]int, 64)
	for i := range cellAges {
		cellAges[i] = make([]int, 64)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(*speed) * time.Millisecond)
	defer ticker.Stop()

	generation := 0

	for {
		select {
		case <-ticker.C:
			// Update cell ages
			for y := 0; y < 64; y++ {
				for x := 0; x < 64; x++ {
					if game.IsAlive(x, y) {
						cellAges[y][x]++
					} else {
						cellAges[y][x] = 0
					}
				}
			}

			// Render to image
			img := pixoo.CreateImage()
			for y := 0; y < 64; y++ {
				for x := 0; x < 64; x++ {
					if game.IsAlive(x, y) {
						c := getCellColor(cellAges[y][x], *colorMode, x, y)
						img.Set(x, y, c)
					} else {
						img.Set(x, y, color.RGBA{0, 0, 0, 255}) // Dead = black
					}
				}
			}

			// Send to display
			if err := client.DrawImage(img); err != nil {
				log.Printf("Error drawing: %v", err)
			}

			// Log stats every 50 generations
			if generation%50 == 0 {
				log.Printf("Generation %d: %d alive cells", generation, game.CountAlive())
			}

			// Step simulation
			game.Step()
			generation++

			// Reseed if population dies out or explodes
			alive := game.CountAlive()
			if alive < 10 {
				log.Println("Population too low, reseeding...")
				game.LoadPattern(*pattern)
				generation = 0
				for i := range cellAges {
					for j := range cellAges[i] {
						cellAges[i][j] = 0
					}
				}
			} else if alive > 3500 {
				log.Println("Population too high, reseeding...")
				game.LoadPattern(*pattern)
				generation = 0
				for i := range cellAges {
					for j := range cellAges[i] {
						cellAges[i][j] = 0
					}
				}
			}

		case <-sigChan:
			log.Println("Shutting down...")
			return
		}
	}
}

func getCellColor(age int, mode string, x, y int) color.RGBA {
	switch mode {
	case "age":
		// Color based on how long the cell has been alive
		if age < 5 {
			return color.RGBA{0, 255, 0, 255} // Young = bright green
		} else if age < 15 {
			return color.RGBA{0, 200, 50, 255} // Medium = green-yellow
		} else if age < 30 {
			return color.RGBA{200, 200, 0, 255} // Older = yellow
		} else {
			return color.RGBA{255, 100, 0, 255} // Ancient = orange
		}

	case "rainbow":
		// Rainbow based on position
		hue := float64((x + y) % 64) / 64.0
		return hueToRGB(hue)

	case "fire":
		// Fire colors
		if age < 3 {
			return color.RGBA{255, 255, 0, 255} // Yellow
		} else if age < 10 {
			return color.RGBA{255, 150, 0, 255} // Orange
		} else {
			return color.RGBA{255, 50, 0, 255} // Red
		}

	case "ocean":
		// Ocean colors
		if age < 5 {
			return color.RGBA{0, 255, 255, 255} // Cyan
		} else if age < 15 {
			return color.RGBA{0, 150, 255, 255} // Light blue
		} else {
			return color.RGBA{0, 50, 200, 255} // Dark blue
		}

	case "matrix":
		// Matrix green
		if age < 3 {
			return color.RGBA{0, 255, 0, 255} // Bright green
		} else if age < 10 {
			return color.RGBA{0, 180, 0, 255} // Medium green
		} else {
			return color.RGBA{0, 100, 0, 255} // Dark green
		}

	default:
		return color.RGBA{0, 255, 0, 255} // Default green
	}
}

func hueToRGB(hue float64) color.RGBA {
	// Convert HSV to RGB (S=1, V=1)
	h := hue * 6.0
	c := 1.0
	x := c * (1.0 - abs(mod(h, 2.0)-1.0))

	var r, g, b float64
	switch int(h) {
	case 0:
		r, g, b = c, x, 0
	case 1:
		r, g, b = x, c, 0
	case 2:
		r, g, b = 0, c, x
	case 3:
		r, g, b = 0, x, c
	case 4:
		r, g, b = x, 0, c
	case 5:
		r, g, b = c, 0, x
	}

	return color.RGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		255,
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func mod(a, b float64) float64 {
	return a - b*float64(int(a/b))
}
