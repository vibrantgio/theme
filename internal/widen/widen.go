// Package widen makes the stems of a TrueType face wider without a change
// to its heights or its spacing.
package widen

import (
	"encoding/binary"
	"errors"
	"math"
	"sort"
)

// Sideways returns a copy of the TrueType font ttf whose simple glyphs are
// wider by em in total. A negative em makes them narrower. Each point
// moves by the x part of its outward miter, by half of em on each side of a
// stroke. No y value changes, and no advance changes. Each glyph's left
// side bearing in hmtx becomes its new xMin. Composite glyphs are copied as
// they are, so they draw their widened components. The copy writes loca in
// the long format.
func Sideways(ttf []byte, em float64) ([]byte, error) {
	be := binary.BigEndian
	if len(ttf) < 12 {
		return nil, errors.New("widen: no table directory")
	}
	numTables := int(be.Uint16(ttf[4:]))
	if len(ttf) < 12+16*numTables {
		return nil, errors.New("widen: short table directory")
	}
	type table struct {
		tag  string
		data []byte
	}
	tables := make([]table, numTables)
	for i := range tables {
		r := ttf[12+16*i:]
		off, n := be.Uint32(r[8:]), be.Uint32(r[12:])
		if uint64(off)+uint64(n) > uint64(len(ttf)) {
			return nil, errors.New("widen: table out of range")
		}
		tables[i] = table{string(r[:4]), ttf[off : off+n]}
	}
	find := func(tag string) int {
		for i := range tables {
			if tables[i].tag == tag {
				return i
			}
		}
		return -1
	}
	hi, li, gi, mi, hhi, hmi := find("head"), find("loca"), find("glyf"), find("maxp"), find("hhea"), find("hmtx")
	if hi < 0 || li < 0 || gi < 0 || mi < 0 || hhi < 0 || hmi < 0 {
		return nil, errors.New("widen: not a TrueType outline font")
	}
	head := append([]byte(nil), tables[hi].data...)
	hmtx := append([]byte(nil), tables[hmi].data...)
	half := em * float64(be.Uint16(head[18:])) / 2
	numGlyphs := int(be.Uint16(tables[mi].data[4:]))
	numMetrics := int(be.Uint16(tables[hhi].data[34:]))
	setLSB := func(g int, v int16) {
		if g < numMetrics {
			be.PutUint16(hmtx[4*g+2:], uint16(v))
		} else {
			be.PutUint16(hmtx[4*numMetrics+2*(g-numMetrics):], uint16(v))
		}
	}
	loca := tables[li].data
	long := be.Uint16(head[50:]) == 1
	at := func(i int) int {
		if long {
			return int(be.Uint32(loca[4*i:]))
		}
		return 2 * int(be.Uint16(loca[2*i:]))
	}
	glyf := tables[gi].data
	var out []byte
	newLoca := make([]byte, 4*(numGlyphs+1))
	for g := 0; g < numGlyphs; g++ {
		be.PutUint32(newLoca[4*g:], uint32(len(out)))
		src := glyf[at(g):at(g+1)]
		if len(src) >= 10 && int16(be.Uint16(src)) > 0 {
			w, err := glyph(src, half)
			if err != nil {
				return nil, err
			}
			src = w
			setLSB(g, int16(be.Uint16(src[2:])))
		}
		out = append(out, src...)
		for len(out)%4 != 0 {
			out = append(out, 0)
		}
	}
	be.PutUint32(newLoca[4*numGlyphs:], uint32(len(out)))
	be.PutUint16(head[50:], 1)
	be.PutUint32(head[8:], 0)
	tables[hi].data, tables[li].data, tables[gi].data, tables[hmi].data = head, newLoca, out, hmtx

	sort.Slice(tables, func(i, j int) bool { return tables[i].tag < tables[j].tag })
	res := append([]byte(nil), ttf[:12]...)
	res = append(res, make([]byte, 16*numTables)...)
	for i, t := range tables {
		for len(res)%4 != 0 {
			res = append(res, 0)
		}
		r := res[12+16*i:]
		copy(r, t.tag)
		be.PutUint32(r[4:], checksum(t.data))
		be.PutUint32(r[8:], uint32(len(res)))
		be.PutUint32(r[12:], uint32(len(t.data)))
		res = append(res, t.data...)
	}
	return res, nil
}

func checksum(b []byte) uint32 {
	var s uint32
	for i := 0; i < len(b); i += 4 {
		var w [4]byte
		copy(w[:], b[i:])
		s += binary.BigEndian.Uint32(w[:])
	}
	return s
}

// errGlyph is the error for a simple glyph whose data ends early.
var errGlyph = errors.New("widen: short glyph")

// glyph rewrites one simple glyph with each x moved by half font units along
// the x part of the point's outward miter. The flags keep only the on-curve
// bit, and every coordinate is written as a two-byte delta.
func glyph(g []byte, half float64) (res []byte, err error) {
	defer func() {
		if recover() != nil {
			res, err = nil, errGlyph
		}
	}()
	be := binary.BigEndian
	nc := int(be.Uint16(g))
	ends := make([]int, nc)
	p := 10
	for i := range ends {
		ends[i] = int(be.Uint16(g[p:]))
		p += 2
	}
	np := ends[nc-1] + 1
	il := int(be.Uint16(g[p:]))
	instr := g[p+2 : p+2+il]
	p += 2 + il
	flags := make([]byte, np)
	for i := 0; i < np; {
		f := g[p]
		p++
		flags[i] = f
		i++
		if f&8 != 0 {
			for r := int(g[p]); r > 0; r-- {
				flags[i] = f
				i++
			}
			p++
		}
	}
	xs, ys := make([]int, np), make([]int, np)
	read := func(v []int, short, same byte) {
		c := 0
		for i, f := range flags {
			switch {
			case f&short != 0:
				d := int(g[p])
				p++
				if f&same == 0 {
					d = -d
				}
				c += d
			case f&same == 0:
				c += int(int16(be.Uint16(g[p:])))
				p += 2
			}
			v[i] = c
		}
	}
	read(xs, 2, 16)
	read(ys, 4, 32)

	// The sign of the summed area gives the winding. TrueType outer
	// contours run clockwise, and their outward normal is (-dy, dx).
	area := 0.0
	start := 0
	for _, e := range ends {
		for i := start; i <= e; i++ {
			j := i + 1
			if j > e {
				j = start
			}
			area += float64(xs[i])*float64(ys[j]) - float64(xs[j])*float64(ys[i])
		}
		start = e + 1
	}
	sgn := 1.0
	if area > 0 {
		sgn = -1
	}
	nx := make([]int, np)
	start = 0
	for _, e := range ends {
		cnt := e - start + 1
		for k := 0; k < cnt; k++ {
			i := start + k
			pr := start + (k-1+cnt)%cnt
			nt := start + (k+1)%cnt
			n1x, n1y := outward(xs[pr], ys[pr], xs[i], ys[i], sgn)
			n2x, n2y := outward(xs[i], ys[i], xs[nt], ys[nt], sgn)
			// The miter is (n1+n2)/(1+n1·n2). The divisor stops at
			// 0.25, so a sharp corner moves at most eight times half.
			d := max(1+n1x*n2x+n1y*n2y, 0.25)
			nx[i] = xs[i] + int(math.Round((n1x+n2x)/d*half))
		}
		start = e + 1
	}
	xmin, xmax := math.MaxInt, math.MinInt
	for _, x := range nx {
		xmin, xmax = min(xmin, x), max(xmax, x)
	}
	res = make([]byte, 0, 12+2*nc+il+5*np)
	res = append(res, g[:10]...)
	be.PutUint16(res[2:], uint16(int16(xmin)))
	be.PutUint16(res[6:], uint16(int16(xmax)))
	res = append(res, g[10:10+2*nc]...)
	res = be.AppendUint16(res, uint16(il))
	res = append(res, instr...)
	for _, f := range flags {
		res = append(res, f&1)
	}
	prev := 0
	for _, x := range nx {
		res = be.AppendUint16(res, uint16(int16(x-prev)))
		prev = x
	}
	prev = 0
	for _, y := range ys {
		res = be.AppendUint16(res, uint16(int16(y-prev)))
		prev = y
	}
	return res, nil
}

// outward is the unit outward normal of the edge from (x0, y0) to (x1, y1),
// or zero for an edge of no length.
func outward(x0, y0, x1, y1 int, sgn float64) (float64, float64) {
	dx, dy := float64(x1-x0), float64(y1-y0)
	l := math.Hypot(dx, dy)
	if l == 0 {
		return 0, 0
	}
	return sgn * -dy / l, sgn * dx / l
}
