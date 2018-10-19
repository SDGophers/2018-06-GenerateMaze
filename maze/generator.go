package maze

import (
	"errors"
	"math/rand"
)

var (
	ErrInvalidMapSize = errors.New("invalid map size")
)

// The smallest map that can be generate and still be valid without any warps
const MinMapArea = 2

// RandMap will generate a new random map with the number of warp points specified.
func RandMap(width, height uint8, warpCount uint8) (*Map, error) {
	miniarea := MinMapArea + warpCount
	if width*height < minarea {
		return nil, ErrInvalidMapSize
	}
	m := Map{
		tiles: make([]tileType, width*height),
		w:     width,
		h:     height,
	}
	// We want to initialize all the tiles to walls.
	for i := range m.tiles {
		m.tiles[i] = Wall
	}
	startX := rand.Intn(int(width))
	startY := rand.Intn(int(height))
	idx := m.idxFor(startX, startY)
	m.tiles[idx] = Start

}

func (m *Map) adjectentIdxes(idx uint8, includeDiagonals bool) []uint8 {

	seedX, seedY := xyFor(idx)

	return []uint8{
		idxFor(seedX+1, seedY+1),
		idxFor(seedX+1, seedY),
		idxFor(seedX, seedY+1),

		idxFor(seedX-1, seedY-1),
		idxFor(seedX-1, seedY),
		idxFor(seedX, seedY-1),

		idxFor(seedX-1, seedY+1),
		idxFor(seedX+1, seedY-1),
	}

}
