# Conway's Game of Life for Pixoo 64

A beautiful implementation of Conway's Game of Life running on your Divoom Pixoo 64!

## Quick Start

Run from **Terminal.app** (not iTerm):

```bash
./game-of-life -host 192.168.1.140
```

## Patterns

Try different starting patterns:

```bash
# Random patterns
./game-of-life -host 192.168.1.140 -pattern random
./game-of-life -host 192.168.1.140 -pattern random-sparse
./game-of-life -host 192.168.1.140 -pattern random-dense

# Classic patterns
./game-of-life -host 192.168.1.140 -pattern gliders
./game-of-life -host 192.168.1.140 -pattern gosper-gun
./game-of-life -host 192.168.1.140 -pattern pulsar
```

## Color Modes

Make it beautiful with different color schemes:

```bash
# Cell age colors (default) - shows how old cells are
./game-of-life -host 192.168.1.140 -color age

# Rainbow gradient
./game-of-life -host 192.168.1.140 -color rainbow

# Fire colors (red/orange/yellow)
./game-of-life -host 192.168.1.140 -color fire

# Ocean colors (blue/cyan)
./game-of-life -host 192.168.1.140 -color ocean

# Matrix green
./game-of-life -host 192.168.1.140 -color matrix
```

## Speed Control

Adjust the simulation speed:

```bash
# Slow (500ms per generation)
./game-of-life -host 192.168.1.140 -speed 500

# Fast (100ms per generation)
./game-of-life -host 192.168.1.140 -speed 100

# Super fast (50ms per generation)
./game-of-life -host 192.168.1.140 -speed 50
```

## Combinations

Combine options for the perfect display:

```bash
# Fast rainbow gliders
./game-of-life -host 192.168.1.140 -pattern gliders -color rainbow -speed 150

# Slow fire Gosper gun
./game-of-life -host 192.168.1.140 -pattern gosper-gun -color fire -speed 300

# Matrix-style random
./game-of-life -host 192.168.1.140 -pattern random -color matrix -speed 100
```

## All Options

```
-host string
    Pixoo 64 device IP address (required)

-pattern string
    Starting pattern (default "random")
    Options: random, random-sparse, random-dense, gliders, gosper-gun, pulsar

-color string
    Color mode (default "age")
    Options: age, rainbow, fire, ocean, matrix

-speed int
    Update speed in milliseconds (default 200)

-brightness int
    Screen brightness 0-100 (default 70)
```

## Features

- **Toroidal topology** - Cells wrap around edges
- **Auto-reseed** - Automatically restarts if population dies or explodes
- **Cell age tracking** - Colors cells based on how long they've been alive
- **Classic patterns** - Includes famous Game of Life patterns
- **64x64 grid** - Perfect for the Pixoo 64 display

## Game of Life Rules

Conway's Game of Life follows these simple rules:

1. Any live cell with 2 or 3 neighbors survives
2. Any dead cell with exactly 3 neighbors becomes alive
3. All other cells die or stay dead

Despite these simple rules, complex and beautiful patterns emerge!

## Recommended Configurations

### Mesmerizing
```bash
./game-of-life -host 192.168.1.140 -pattern random -color age -speed 150
```

### Classic
```bash
./game-of-life -host 192.168.1.140 -pattern gosper-gun -color matrix -speed 200
```

### Psychedelic
```bash
./game-of-life -host 192.168.1.140 -pattern random-dense -color rainbow -speed 100
```

### Chill
```bash
./game-of-life -host 192.168.1.140 -pattern pulsar -color ocean -speed 400
```

## Building

To rebuild after changes:

```bash
go build -ldflags="-linkmode=external" -o game-of-life game-of-life.go
```

## Troubleshooting

If you get "no route to host" errors:
1. Make sure you're running from **Terminal.app** (not iTerm)
2. The first time you run it, macOS should prompt for Local Network access - click "Allow"
3. See the main README.md for more troubleshooting

Enjoy watching life evolve on your Pixoo 64! 🎨✨
