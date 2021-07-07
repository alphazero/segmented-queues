// Doost

package segque

import "fmt"

/// Container types ////////////////////////////////////////////
type CType int

func (c CType) String() string {
	return ctypes[c]
}

const (
	BA       CType = iota // basic array with direct addressing
	Co2_I_C               // single array choice of 2 using container sequence number
	Co2_I_R               // single array choice of 2 using record sequence number
	Co2_II_C              // double array choice of 2 using container sequence number
	Co2_II_R              // double array choice of 2 using record sequence number
	Co4_IV_C              // quad array choice of 2 using container sequence number
	Co4_IV_R              // quad array choice of 2 using record sequence number
)

var ctypes = map[CType]string{
	BA:       "BA",
	Co2_I_C:  "Co2_I_C",
	Co2_I_R:  "Co2_I_R",
	Co2_II_C: "Co2_II_C",
	Co2_II_R: "Co2_II_R",
	Co4_IV_C: "Co4_IV_C",
	Co4_IV_R: "Co4_IV_R",
}

/// container container ///////////////////////////////////////////////
type container struct {
	ctype    CType // container type
	capacity int
	mask     uint64 // mask used to assign key to bucket
	seqmask  uint64 // sequence number mask to emulate rollover
	arrcnt   int
	arr      [][]*FifoQ // backing arrays
	seqnum   []uint64   // one per backing array
}

func (p *container) String() string {
	return fmt.Sprintf("type:%s mask:%x seqmask:%x capacity:%d", p.ctype, p.mask, p.seqmask, p.capacity)
}

func (p *container) ArrCnt() int { return p.arrcnt }

// Update supports Container.Update
// REVU use PickOldest
func (c *container) Update(p *Params, seqnum uint64, key ...uint64) uint64 {
	var idxs = make([]int, c.arrcnt)
	for i := 0; i < c.arrcnt; i++ {
		idxs[i] = int(key[i] & c.mask)
	}
	// debug
	for i := 0; i < c.arrcnt; i++ {
		Trace(p, "idx%d: %d => ", i, idxs[i])
		c.arr[i][idxs[i]].DebugPrint(p)
	}
	var pick = 0
	var seqnums = make([]uint64, c.arrcnt)
	switch c.ctype {
	case BA:
		// NOP - pick is 0 so we are picking idxs[0] as required
	case Co2_I_C, Co2_II_C, Co4_IV_C:
		for i := 0; i < c.arrcnt; i++ {
			seqnums[i] = c.arr[i][idxs[i]].Seqnum()
		}
		pick = PickOldest(p, c.seqmask, seqnum, seqnums)
	case Co2_I_R, Co2_II_R, Co4_IV_R:
		for i := 0; i < c.arrcnt; i++ {
			seqnums[i] = c.arr[i][idxs[i]].Tail()
		}
		pick = PickOldest(p, c.seqmask, seqnum, seqnums)
	}
	idx := idxs[pick]
	arr := c.arr[pick]
	ev := arr[idx].Add(seqnum)
	Trace(p, "evict %x => ", ev)
	arr[idx].DebugPrint(p)
	Trace(p, "------------\n")
	return ev
	//	return arr[idx].Add(seqnum)
}

/// public api ///////////////////////////////////////////////////

// Container defines api for updating a container
type Container interface {
	ArrCnt() int
	// op sequence number is the full bits sequence number.
	// keys 1 or more are used for selecting container bucket
	// returns evicted seqnum - 0 is zero value
	Update(p *Params, seqnum uint64, key ...uint64) uint64
	// container descriptive string
	String() string
}

// NewContainer creates a new container of specified CType, with allocated FifoQ array(s)
// as required. Buckets are evenly divided across two arrays for the double array types,
// with key mask adjusted accordingly.
func NewContainer(ctype CType, buckets int, slots int, seqmask uint64) Container {

	var arrcnt int
	switch ctype {
	case BA, Co2_I_C, Co2_I_R:
		arrcnt = 1
	case Co2_II_C, Co2_II_R:
		arrcnt = 2
	case Co4_IV_C, Co4_IV_R:
		arrcnt = 4
	}
	var arrlen = buckets / arrcnt
	var mask = uint64(arrlen - 1)
	c := &container{
		ctype:    ctype,
		capacity: buckets * slots,
		mask:     mask,
		seqmask:  seqmask,
		arr:      make([][]*FifoQ, arrcnt),
		seqnum:   make([]uint64, arrcnt),
		arrcnt:   arrcnt,
	}
	for i := 0; i < arrcnt; i++ {
		c.arr[i] = make([]*FifoQ, arrlen)
		for j := 0; j < arrlen; j++ {
			c.arr[i][j] = NewFifoQ(slots)
		}
	}
	return c
}

/// sequence number algorithm

// pickOldesst applies the seqmask to emulate a mask sized roll-over counter.
// it also applies the same mask to the array of sequnece numbers provided.
// the algorithm assumes that any given sequence number is at most one cycle
// behind the counter. thus, if any of the sequence numbers in the array are
// greater than the (masked) seqnum, it is assumed they are lagging a cycle
// behind.
func PickOldest(p *Params, seqmask uint64, seqnum uint64, seqnums []uint64) (pickIdx int) {
	var refnum = seqmask & seqnum
	var least uint64 = seqmask << 2
	var cycle = seqmask + 1
	var idx int
	for i, v := range seqnums {
		v0 := v & seqmask
		// if v0 exceed refnum then it is a value form previous cycle.
		// so if less add a cycle value to it - this was
		// we can reasonably compare two seqnums to see which is less
		if v0 < refnum {
			v0 += cycle
		}
		if v0 < least {
			least = v0
			idx = i
		}
	}
	Trace(p, "with mask %x seqnum %x refnum %x pick from %x - picked %x at idx %d\n", seqmask, seqnum, refnum, seqnums, seqnums[idx], idx)
	return idx
}
