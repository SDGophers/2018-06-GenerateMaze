package maze

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

type tileType uint8

const (
	headerBom = iota
	headerMagic1
	headerMagic2
	headerVersion
	headerWidth
	headerHeight
	headerLen
)

const (
	MagicNumber1 = 0x5d
	MagicNumber2 = 0x90
	FileVersion  = 1
)

const (
	path = tileType(iota)
	wall
	start
	rgoal
	ogoal
	warp
)

const (
	itemRequiredGoal = 0
	itemOptionalGoal = 1
	itemWarp         = 2
)

type item struct {
	Typ  byte
	X, Y byte
}

type Map struct {
	w, h  byte
	tiles []tileType
	start [2]byte
	warps map[int]int
}

func (m *Map) idxFor(x, y byte) int { return (int(y) * int(m.w)) + int(x) }

func ReadMap(r io.Reader) (*Map, error) {

	var bytes = make([]byte, headerLen)
	var bom binary.ByteOrder
	var err error
	// the bom is the first byte of the file.
	if _, err = r.Read(bytes); err != nil {
		return nil, err
	}

	switch bytes[headerBom] {
	case 0:
		bom = binary.BigEndian
	case 1:
		bom = binary.LittleEndian
	default:
		return nil, errors.New("bad byte order marker")
	}

	// Magic and Version
	if bytes[headerMagic1] != MagicNumber1 || bytes[headerMagic2] != MagicNumber2 || bytes[headerVersion] != FileVersion {
		return nil, errors.New("Bad magic or version number.")
	}
	wl, hl := bytes[headerWidth], bytes[headerHeight]
	if (int(wl) * int(hl)) <= 1 {
		return nil, errors.New("invalid map, not big enough.")
	}

	// Lets read the values of the map
	var n int16
	if err = binary.Read(r, bom, &n); err != nil {
		return nil, err
	}
	if int(wl)*int(hl) > int(n)*8 {
		return nil, errors.New("invalid map, not enough map data.")
	}

	var mdata = make([]byte, int(n))
	if err = binary.Read(r, bom, &mdata); err != nil {
		return nil, err
	}

	var m Map
	m.w = wl
	m.h = hl
	m.warps = make(map[int]int)
	// height and width
	w, h := byte(0), byte(0)
	m.tiles = make([]tileType, int(wl)*int(hl))
	for _, b := range mdata {
		// each byte contains 8 items.
		for j := uint(0); j < 8; j++ {
			if (b<<j)&0x80 != 0 {
				m.tiles[m.idxFor(w, h)] = wall
			}
			w++
			if w >= wl {
				w, h = 0, h+1
			}
		}
	}
	if err = binary.Read(r, bom, &m.start); err != nil {
		return nil, err
	}
	if m.tiles[m.idxFor(m.start[0], m.start[1])] != path {
		return nil, errors.New("invalid start location.")
	}
	m.tiles[m.idxFor(m.start[0], m.start[1])] = start

	var numberItems byte
	if err = binary.Read(r, bom, &numberItems); err != nil {
		return nil, err
	}
	var hasRequiredGoal bool
	var it item

	for i := 0; i < int(numberItems); i++ {

		if err = binary.Read(r, bom, &it); err != nil {
			return nil, err
		}

		idx := m.idxFor(it.X, it.Y)

		if m.tiles[idx] != path {
			return nil, errors.New("invalid item location.")
		}

		switch it.Typ {
		case itemRequiredGoal:
			hasRequiredGoal = true
			m.tiles[idx] = rgoal
		case itemOptionalGoal:
			m.tiles[idx] = ogoal

		case itemWarp:
			m.tiles[idx] = warp
			to := make([]byte, 2)
			if err = binary.Read(r, bom, &to); err != nil {
				return nil, err
			}
			m.warps[idx] = m.idxFor(to[0], to[1])
		}
	}
	if !hasRequiredGoal {
		return nil, errors.New("At least one required goal is required.")
	}
	return &m, nil
}

func (m *Map) String() string {
	if m == nil {
		return ""
	}
	var val strings.Builder
	val.Grow(256)

	for y := byte(0); y < m.h; y++ {
		for x := byte(0); x < m.w; x++ {
			idx := m.idxFor(x, y)
			switch m.tiles[idx] {
			case path:
				val.WriteRune(' ')
				//val.WriteRune(' ')
			case wall:
				val.WriteRune('█')
				//val.WriteRune('█')
			case start:
				//val.WriteRune('🦆')
				val.WriteRune('S')
			case rgoal:
				//val.WriteRune('🐣')
				val.WriteRune('G')
			case ogoal:
				//val.WriteRune('🐥')
				val.WriteRune('O')
			case warp:
				//val.WriteRune('🕳')
				val.WriteRune('W')
			}

		}
		val.WriteRune('\n')
	}
	return val.String()
}
