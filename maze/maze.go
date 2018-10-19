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
	MSb          = 0x80 // Most sigficent bit -- 1000_0000

)

const (
	Path = tileType(iota)
	Wall
	Start
	Rgoal
	Ogoal
	Warp
)

const (
	itemRequiredGoal = 0
	itemOptionalGoal = 1
	itemWarp         = 2
)

var (
	ErrInvalidMapNotLargeEnough = errors.New("invalid map, not big enough")
	ErrInvalidMapNotEnoughData  = errors.New("invalid map, not enough data")
	ErrInvalidMapInvalidStart   = errors.New("invalid map, start location")
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
func (m *Map) xyFor(idx int) (byte, byte) {
	x := idx % m.w
	y := idx / m.w
}

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

	areaInBits := int(wl) * int(hl)

	if areaInBits <= 1 {
		return nil, ErrInvalidMapNotLargeEnough
	}

	// Lets read the values of the map
	var n int16
	if err = binary.Read(r, bom, &n); err != nil {
		return nil, err
	}

	dataInBits := int(n) * 8

	if areaInBits > dataInBits {
		return nil, ErrInvalidMapNotEnoughData
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

			// Check to see if the the most sigificent bit is set.
			if (b<<j)&MSb != 0 {
				m.tiles[m.idxFor(w, h)] = Wall
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
	if m.tiles[m.idxFor(m.start[0], m.start[1])] != Path {
		return nil, ErrInvalidMapInvalidStart
	}
	m.tiles[m.idxFor(m.start[0], m.start[1])] = Start

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

		if m.tiles[idx] != Path {
			return nil, errors.New("invalid item location.")
		}

		switch it.Typ {
		case itemRequiredGoal:
			hasRequiredGoal = true
			m.tiles[idx] = Rgoal
		case itemOptionalGoal:
			m.tiles[idx] = Ogoal

		case itemWarp:
			m.tiles[idx] = Warp
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
			case Path:
				val.WriteRune(' ')
			case Wall:
				val.WriteRune('█')
			case Start:
				val.WriteRune('S')
			case Rgoal:
				val.WriteRune('G')
			case Ogoal:
				val.WriteRune('O')
			case Warp:
				val.WriteRune('W')
			}

		}
		val.WriteRune('\n')
	}
	return val.String()
}
