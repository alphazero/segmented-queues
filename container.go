// Doost

package segque

import "fmt"

type Container interface {
	// op sequence number
	// keys 1 or more are used for selecting container bucket
	// returns evicted seqnum - 0 is zero value
	Update(seqnum uint64, key ...uint64) uint64
}

/// basic container //////////////////////////////////////////////
type basic_c struct {
	mask uint64
	arr  []*FifoQ
}

//func NewBasicContainer(buckets int, slots int, seqmask uint64) *basic_c {
func NewBasicContainer(buckets int, slots int) *basic_c {
	c := basic_c{
		mask: uint64(buckets - 1),
		arr:  make([]*FifoQ, buckets),
	}
	for i := 0; i < len(c.arr); i++ {
		c.arr[i] = NewFifoQ(slots)
	}
	return &c
}

// Basic Array only supports direct addressing
func (c *basic_c) Update(seqnum uint64, key ...uint64) uint64 {
	idx := key[0] & c.mask
	return c.arr[idx].Add(seqnum)
}

/// Co2 container base ///////////////////////////////////////////
type co2_base struct {
	mask    uint64
	seqmask uint64
	ctype   CType
	//	cseq    bool
}

/// Co2 I container //////////////////////////////////////////////

// supports for co2_i and co2_i_r
type co2_i_c struct {
	base co2_base
	arr  []*FifoQ
}

//func NewCo2IContainer(buckets int, slots int, seqmask uint64, cseq bool) *co2_i_c {
func NewCo2IContainer(buckets int, slots int, seqmask uint64, ctype CType) *co2_i_c {
	switch ctype {
	case Co2_I_C, Co2_I_R:
	default:
		panic(fmt.Sprintf("invalid strategy %d %s for single array container", ctype, Strategies[ctype]))
	}
	c := co2_i_c{
		base: co2_base{
			mask:    uint64(buckets - 1),
			seqmask: seqmask,
			ctype:   ctype,
			//cseq:    cseq,
		},
		arr: make([]*FifoQ, buckets),
	}
	for i := 0; i < len(c.arr); i++ {
		c.arr[i] = NewFifoQ(slots)
	}
	return &c
}

// Single backing array and choice-of-two strategies
func (c *co2_i_c) Update(seqnum uint64, key ...uint64) uint64 {
	idx0 := int(key[0] & c.base.mask)
	idx1 := int(key[1] & c.base.mask)
	var choice = idx0
	switch c.base.ctype {
	case Co2_I_C:
		if (c.arr[idx0].Seqnum() & c.base.seqmask) > (c.arr[idx1].Seqnum() & c.base.seqmask) {
			choice = idx1
		}
	case Co2_I_R:
		if (c.arr[idx0].Tail() & c.base.seqmask) > (c.arr[idx1].Tail() & c.base.seqmask) {
			choice = idx1
		}
	}

	return c.arr[choice].Add(seqnum)
}

/// Co2 II container /////////////////////////////////////////////

// supports for co2_ii and co2_ii_r
type co2_ii_c struct {
	base co2_base
	arr0 []*FifoQ
	arr1 []*FifoQ
}

//func NewCo2IIContainer(buckets int, slots int, seqmask uint64, cseq bool) *co2_ii_c {
func NewCo2IIContainer(buckets int, slots int, seqmask uint64, ctype CType) *co2_ii_c {
	c := co2_ii_c{
		base: co2_base{
			mask:    uint64(buckets-1) / 2,
			seqmask: seqmask,
			ctype:   ctype,
			//cseq:    cseq,
		},
		arr0: make([]*FifoQ, buckets/2),
		arr1: make([]*FifoQ, buckets/2),
	}
	for i := 0; i < len(c.arr0); i++ {
		c.arr0[i] = NewFifoQ(slots)
	}
	for i := 0; i < len(c.arr1); i++ {
		c.arr1[i] = NewFifoQ(slots)
	}
	return &c
}

// Container with two backing arrays with choice-of-two strategies including random pick.
func (c *co2_ii_c) Update(seqnum uint64, key ...uint64) uint64 {
	idx0 := int(key[0] & c.base.mask)
	idx1 := int(key[1] & c.base.mask)
	var arr = c.arr0
	var choice = idx0
	switch c.base.ctype {
	case Co2_II_C:
		if (c.arr0[idx0].Seqnum() & c.base.seqmask) > (c.arr1[idx1].Seqnum() & c.base.seqmask) {
			choice = idx1
			arr = c.arr1
		}
	case Co2_II_R:
		if (c.arr0[idx0].Tail() & c.base.seqmask) > (c.arr1[idx1].Tail() & c.base.seqmask) {
			choice = idx1
			arr = c.arr1
		}
	case Co2_II_Rand:
		if 0x8000000000000000&key[1] == 0x8000000000000000 {
			choice = idx1
			arr = c.arr1
		}
	}

	return arr[choice].Add(seqnum)
}
