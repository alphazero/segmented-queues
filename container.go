// Doost

package segque

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
	Co2_II_Rand       // double array choice of 2 with random choice
)

var ctypes = map[CType]string{
	BA:          "BA",
	Co2_I_C:     "Co2_I_C",
	Co2_I_R:     "Co2_I_R",
	Co2_II_C:    "Co2_II_C",
	Co2_II_R:    "Co2_II_R",
	Co2_II_Rand: "Co2_II_Rand",
}

/// interface ////////////////////////////////////////////////////
type Container interface {
	// op sequence number is the full bits sequence number.
	// keys 1 or more are used for selecting container bucket
	// returns evicted seqnum - 0 is zero value
	Update(seqnum uint64, key ...uint64) uint64
}

/// basic container //////////////////////////////////////////////
type base struct {
	mask    uint64
	seqmask uint64
	ctype   CType
}

// container with one backing array
type one_barr struct {
	base
	arr []*FifoQ
}

func (c *one_barr) Update(seqnum uint64, key ...uint64) uint64 {
	var idx int
	switch c.base.ctype {
	case BA:
		// basic addressing with mask
		idx = int(key[0] & c.mask)
	case Co2_I_C:
		// pick lower container sequence number
		idx0 := int(key[0] & c.mask)
		idx1 := int(key[1] & c.mask)
		idx = idx0
		if (c.arr[idx0].Seqnum() & c.seqmask) > (c.arr[idx1].Seqnum() & c.seqmask) {
			idx = idx1
		}
	case Co2_I_R:
		// pick lower record sequence number
		idx0 := int(key[0] & c.mask)
		idx1 := int(key[1] & c.mask)
		idx = idx0
		if (c.arr[idx0].Tail() & c.seqmask) > (c.arr[idx1].Tail() & c.seqmask) {
			idx = idx1
		}
	}
	return c.arr[idx].Add(seqnum)
}

// container with two backing arrays
type two_barr struct {
	base
	arr1 []*FifoQ
	arr2 []*FifoQ
}

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
