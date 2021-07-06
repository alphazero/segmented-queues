// Doost

package segque

import "fmt"

/// Container types ////////////////////////////////////////////
type CType int

func (c CType) String() string {
	return ctypes[c]
}

const (
	_           CType = iota
	BA                // basic array with direct addressing
	Co2_I_C           // single array choice of 2 using container sequence number
	Co2_I_R           // single array choice of 2 using record sequence number
	Co2_II_C          // double array choice of 2 using container sequence number
	Co2_II_R          // double array choice of 2 using record sequence number
	Co2_II_Rand       // double array choice of 2 with random choice // REVU get rid of rand
	Co3_III_C         // triple array choice of 2 using container sequence number
	Co3_III_R         // triple array choice of 2 using record sequence number
)

var ctypes = map[CType]string{
	BA:          "BA",
	Co2_I_C:     "Co2_I_C",
	Co2_I_R:     "Co2_I_R",
	Co2_II_C:    "Co2_II_C",
	Co2_II_R:    "Co2_II_R",
	Co2_II_Rand: "Co2_II_Rand",
	Co3_III_C:   "Co3_III_C",
	Co3_III_R:   "Co3_III_R",
}

/// container base ///////////////////////////////////////////////
type base struct {
	mask    uint64 // mask used to assign key to bucket
	seqmask uint64 // sequence number mask to emulate rollover
	ctype   CType  // container type
}

func (p *base) String() string {
	return fmt.Sprintf("type:%s mask:%x seqmask:%x", p.ctype, p.mask, p.seqmask)
}

// type I container - container with one backing array
type one_barr struct {
	base
	arr    []*FifoQ
	seqnum uint64 // REVU use these later when testing seqnum bit lengths
}

func (c *one_barr) String() string {
	return fmt.Sprintf("%s size:%d", c.base.String(), len(c.arr))
}

// Update supports Container.Update
// REVU TODO use PickOldest
func (c *one_barr) Update(seqnum uint64, key ...uint64) uint64 {
	var idx int
	switch c.base.ctype {
	case BA:
		// basic addressing with mask
		idx = int(key[0] & c.mask)
	case Co2_I_C:
		// pick lower container sequence number
		//		idx0 := int(key[0] & c.mask)
		//		idx1 := int(key[1] & c.mask)
		//		idx = idx0
		var idxs = []int{int(key[0] & c.mask), int(key[1] & c.mask)}
		var seqnums = []uint64{c.arr[idxs[0]].Seqnum(), c.arr[idxs[1]].Seqnum()}
		pick := PickOldest(c.seqmask, seqnum, seqnums)
		idx = idxs[pick]
		//		if (c.arr[idx0].Seqnum() & c.seqmask) > (c.arr[idx1].Seqnum() & c.seqmask) {
		//			idx = idx1
		//		}
	case Co2_I_R:
		// pick lower record sequence number
		//		idx0 := int(key[0] & c.mask)
		//		idx1 := int(key[1] & c.mask)
		//		idx = idx0
		var idxs = []int{int(key[0] & c.mask), int(key[1] & c.mask)}
		var seqnums = []uint64{c.arr[idxs[0]].Tail(), c.arr[idxs[1]].Tail()}
		pick := PickOldest(c.seqmask, seqnum, seqnums)
		idx = idxs[pick]
		//		if (c.arr[idx0].Tail() & c.seqmask) > (c.arr[idx1].Tail() & c.seqmask) {
		//			idx = idx1
		//		}
	}
	return c.arr[idx].Add(seqnum)
}

// type II container - container with two backing arrays
type two_barr struct {
	base
	arr1    []*FifoQ
	arr2    []*FifoQ
	seqnum1 uint64
	seqnum2 uint64
}

func (c *two_barr) String() string {
	return fmt.Sprintf("%s size:%d (x2)", c.base.String(), len(c.arr1))
}

// Update supports Container.Update
// REVU TODO use PickOldest
func (c *two_barr) Update(seqnum uint64, key ...uint64) uint64 {
	idx1 := int(key[0] & c.mask)
	idx2 := int(key[1] & c.mask)
	var idx = idx1   // initial choice
	var arr = c.arr1 // initial choice
	switch c.base.ctype {
	case Co2_II_Rand:
		// use hi bit to flip a coin
		if 0x8000000000000000&key[1] == 0x8000000000000000 {
			idx = idx2
			arr = c.arr2
		}
	case Co2_II_C:
		// pick lower container sequence number
		var idxs = []int{int(key[0] & c.mask), int(key[1] & c.mask)}
		var seqnums = []uint64{c.arr1[idxs[0]].Tail(), c.arr2[idxs[1]].Tail()}
		pick := PickOldest(c.seqmask, seqnum, seqnums)
		idx = idxs[pick]
		if (c.arr1[idx1].Seqnum() & c.seqmask) > (c.arr2[idx2].Seqnum() & c.seqmask) {
			idx = idx2
			arr = c.arr2
		}
	case Co2_II_R:
		// pick lower record sequence number
		if (c.arr1[idx1].Tail() & c.seqmask) > (c.arr2[idx2].Tail() & c.seqmask) {
			idx = idx2
			arr = c.arr2
		}
	}
	return arr[idx].Add(seqnum)
}

/// public api ///////////////////////////////////////////////////

// Container defines api for updating a container
type Container interface {
	// op sequence number is the full bits sequence number.
	// keys 1 or more are used for selecting container bucket
	// returns evicted seqnum - 0 is zero value
	Update(seqnum uint64, key ...uint64) uint64
	// container descriptive string
	String() string
}

// NewContainer creates a new container of specified CType, with allocated FifoQ array(s)
// as required. Buckets are evenly divided across two arrays for the double array types,
// with key mask adjusted accordingly.
func NewContainer(ctype CType, buckets int, slots int, seqmask uint64) Container {
	var container Container

	switch ctype {
	case BA, Co2_I_C, Co2_I_R:
		mask := uint64(buckets - 1)
		c := one_barr{
			base: base{
				mask:    mask,
				seqmask: seqmask,
				ctype:   ctype,
			},
			arr: make([]*FifoQ, buckets),
		}
		for i := 0; i < len(c.arr); i++ {
			c.arr[i] = NewFifoQ(slots)
		}
		container = &c
	case Co2_II_C, Co2_II_R, Co2_II_Rand:
		size := buckets / 2
		mask := uint64(size - 1)
		c := two_barr{
			base: base{
				mask:    mask,
				seqmask: seqmask,
				ctype:   ctype,
			},
			arr1: make([]*FifoQ, size),
			arr2: make([]*FifoQ, size),
		}
		for i := 0; i < size; i++ {
			c.arr1[i] = NewFifoQ(slots)
			c.arr2[i] = NewFifoQ(slots)
		}
		container = &c
	}
	return container
}

/// sequence number algorithm

// pickOldesst applies the seqmask to emulate a mask sized roll-over counter.
// it also applies the same mask to the array of sequnece numbers provided.
// the algorithm assumes that any given sequence number is at most one cycle
// behind the counter. thus, if any of the sequence numbers in the array are
// greater than the (masked) seqnum, it is assumed they are lagging a cycle
// behind.
func PickOldest(seqmask uint64, seqnum uint64, seqnums []uint64) (pickIdx int) {
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
	return idx
}
