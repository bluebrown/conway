package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"
)

type field [height][width]bool

// calculate the next generation
func (f *field) step() field {

	next := field{}

	var (
		rowTop int
		rowBot int
		colLef int
		colRig int
	)

	// iterate cells in game grid
	for rowIdx, row := range f {

		// find top and bottom row
		// wraps around if index is out of bounds
		rowTop = rowIdx - 1
		if rowTop < 0 {
			rowTop = len(f) - 1
		}

		rowBot = rowIdx + 1
		if rowBot > len(f)-1 {
			rowBot = 0
		}

		for colIdx, cell := range row {

			// find left and right column
			// wraps around if index is out of bounds
			colLef = colIdx - 1
			if colLef < 0 {
				colLef = len(row) - 1
			}

			colRig = colIdx + 1
			if colRig > len(row)-1 {
				colRig = 0
			}

			// count 8 neighbors
			// continue early if more than 3 neighers are found
			ngbrs := 0
			// top left
			if f[rowTop][colLef] {
				ngbrs++
			}
			// top
			if f[rowTop][colIdx] {
				ngbrs++
			}
			// top right
			if f[rowTop][colRig] {
				ngbrs++
			}
			// left
			if f[rowIdx][colLef] {
				ngbrs++
			}
			// right
			if f[rowIdx][colRig] {
				ngbrs++
			}
			// bottom left
			if f[rowBot][colLef] {
				ngbrs++
			}
			// bottom
			if f[rowBot][colIdx] {
				ngbrs++
			}
			// bottom right
			if f[rowBot][colRig] {
				ngbrs++
			}

			// conways rules of life
			if cell && ngbrs == 2 || ngbrs == 3 {
				next[rowIdx][colIdx] = true
				continue
			}

		}
	}
	*f = next
	return *f
}

// marshal to json encoded byte array,
// used for transfering the data
func (f field) encode() []byte {
	b, _ := json.Marshal(f)
	return b
}

// randomize the field
func (f *field) random() field {
	rand.Seed(time.Now().UnixNano())
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			f[y][x] = rand.Intn(2) == 1
		}
	}
	return *f
}

// call step every game tick until the context is done
func (f *field) run(ctx context.Context, results chan []byte) {
	results <- f.encode()
	f.step()
	for {
		select {
		case <-ctx.Done():
			close(results)
			return
		case <-time.After(tickspeed):
			results <- f.encode()
			f.step()
		}
	}
}

// start new game and return channel to receive computed steps
func newGame(ctx context.Context) chan []byte {
	results := make(chan []byte)
	cw := field{}
	cw.random()
	go cw.run(ctx, results)
	return results
}
