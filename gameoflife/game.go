package gameoflife

import (
	"math/rand"
)

const Size = 64

type Game struct {
	cells     [Size][Size]bool
	nextCells [Size][Size]bool
	generation int
}

func NewGame() *Game {
	return &Game{}
}

// RandomSeed initializes the grid with random alive cells
func (g *Game) RandomSeed(density float64) {
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			g.cells[y][x] = rand.Float64() < density
		}
	}
	g.generation = 0
}

// Glider creates a glider pattern
func (g *Game) Glider(x, y int) {
	pattern := [][]int{
		{1, 0},
		{2, 1},
		{0, 2}, {1, 2}, {2, 2},
	}
	for _, p := range pattern {
		px, py := x+p[0], y+p[1]
		if px >= 0 && px < Size && py >= 0 && py < Size {
			g.cells[py][px] = true
		}
	}
}

// Patterns creates various interesting starting patterns
func (g *Game) LoadPattern(name string) {
	g.Clear()

	switch name {
	case "gliders":
		g.Glider(5, 5)
		g.Glider(15, 20)
		g.Glider(30, 10)
		g.Glider(45, 35)

	case "gosper-gun":
		// Left square
		g.cells[5][1] = true
		g.cells[5][2] = true
		g.cells[6][1] = true
		g.cells[6][2] = true

		// Left part
		g.cells[5][11] = true
		g.cells[6][11] = true
		g.cells[7][11] = true
		g.cells[4][12] = true
		g.cells[8][12] = true
		g.cells[3][13] = true
		g.cells[9][13] = true
		g.cells[3][14] = true
		g.cells[9][14] = true
		g.cells[6][15] = true
		g.cells[4][16] = true
		g.cells[8][16] = true
		g.cells[5][17] = true
		g.cells[6][17] = true
		g.cells[7][17] = true
		g.cells[6][18] = true

		// Middle part
		g.cells[3][21] = true
		g.cells[4][21] = true
		g.cells[5][21] = true
		g.cells[3][22] = true
		g.cells[4][22] = true
		g.cells[5][22] = true
		g.cells[2][23] = true
		g.cells[6][23] = true
		g.cells[1][25] = true
		g.cells[2][25] = true
		g.cells[6][25] = true
		g.cells[7][25] = true

		// Right square
		g.cells[3][35] = true
		g.cells[4][35] = true
		g.cells[3][36] = true
		g.cells[4][36] = true

	case "random":
		g.RandomSeed(0.3)

	case "random-sparse":
		g.RandomSeed(0.15)

	case "random-dense":
		g.RandomSeed(0.5)

	case "pulsar":
		// Pulsar oscillator
		pattern := [][]int{
			// Top section
			{2, 0}, {3, 0}, {4, 0}, {8, 0}, {9, 0}, {10, 0},
			{0, 2}, {5, 2}, {7, 2}, {12, 2},
			{0, 3}, {5, 3}, {7, 3}, {12, 3},
			{0, 4}, {5, 4}, {7, 4}, {12, 4},
			{2, 5}, {3, 5}, {4, 5}, {8, 5}, {9, 5}, {10, 5},
			// Middle gap at y=6
			{2, 7}, {3, 7}, {4, 7}, {8, 7}, {9, 7}, {10, 7},
			{0, 8}, {5, 8}, {7, 8}, {12, 8},
			{0, 9}, {5, 9}, {7, 9}, {12, 9},
			{0, 10}, {5, 10}, {7, 10}, {12, 10},
			{2, 12}, {3, 12}, {4, 12}, {8, 12}, {9, 12}, {10, 12},
		}
		centerX, centerY := Size/2-6, Size/2-6
		for _, p := range pattern {
			x, y := centerX+p[0], centerY+p[1]
			if x >= 0 && x < Size && y >= 0 && y < Size {
				g.cells[y][x] = true
			}
		}
	}

	g.generation = 0
}

// Clear resets all cells to dead
func (g *Game) Clear() {
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			g.cells[y][x] = false
		}
	}
	g.generation = 0
}

// CountNeighbors counts alive neighbors for a cell
func (g *Game) CountNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			// Wrap around edges (toroidal topology)
			nx := (x + dx + Size) % Size
			ny := (y + dy + Size) % Size

			if g.cells[ny][nx] {
				count++
			}
		}
	}
	return count
}

// Step advances the game by one generation
func (g *Game) Step() {
	// Calculate next generation
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			neighbors := g.CountNeighbors(x, y)
			alive := g.cells[y][x]

			// Conway's Game of Life rules:
			// 1. Any live cell with 2 or 3 neighbors survives
			// 2. Any dead cell with exactly 3 neighbors becomes alive
			// 3. All other cells die or stay dead
			if alive {
				g.nextCells[y][x] = neighbors == 2 || neighbors == 3
			} else {
				g.nextCells[y][x] = neighbors == 3
			}
		}
	}

	// Swap buffers
	g.cells = g.nextCells
	g.generation++
}

// IsAlive returns whether a cell is alive
func (g *Game) IsAlive(x, y int) bool {
	if x < 0 || x >= Size || y < 0 || y >= Size {
		return false
	}
	return g.cells[y][x]
}

// Generation returns the current generation number
func (g *Game) Generation() int {
	return g.generation
}

// CountAlive returns the total number of alive cells
func (g *Game) CountAlive() int {
	count := 0
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			if g.cells[y][x] {
				count++
			}
		}
	}
	return count
}
