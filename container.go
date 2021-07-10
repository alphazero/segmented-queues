// Doost

package segque

import "fmt"

/// Container //////////////////////////////////////////////////

// Container defines api for updating a container
type Container interface {
	// Returns the number of backing arrays
	ArrCnt() int
	// Returns the choice-of-N arity
	Choices() int
	// Returns true if choice uses container sequence num
	UseCSeqnum() bool
	// op sequence number is the full bits sequence number.
	// keys 1 or more are used for selecting container bucket
	// returns evicted seqnum - 0 is zero value
	Update(p *Params, seqnum uint64, key ...uint64) uint64
	// container descriptive string
	String() string
}

/// Container types ////////////////////////////////////////////

type CType int

func (c CType) String() string {
	return ctypes[c]
}

const (
	BA      CType = iota // basic array with direct addressing
	C2_A1_C              // single array choice of 2 using container sequence number
	C2_A1_R              // single array choice of 2 using record sequence number
	C2_A2_C              // double array choice of 2 using container sequence number
	C2_A2_R              // double array choice of 2 using record sequence number
	C4_A4_C              // quad array choice of 4 using container sequence number
	C4_A4_R              // quad array choice of 4 using record sequence number
	C8_A8_C              // octal array choice of 8 using container sequence number
	C8_A8_R              // octal array choice of 8 using record sequence number
)

var ctypes = map[CType]string{
	BA:      "BA",
	C2_A1_C: "C2_A1_C",
	C2_A1_R: "C2_A1_R",
	C2_A2_C: "C2_A2_C",
	C2_A2_R: "C2_A2_R",
	C4_A4_C: "C4_A4_C",
	C4_A4_R: "C4_A4_R",
	C8_A8_C: "C8_A8_C",
	C8_A8_R: "C8_A8_R",
}

/// Container support /////////////////////////////////////////////////

// container Container interface for various container configurations
// and load balancing policies. Capacity, backing array count and n-arry
// choice policies are all power of 2 based.
type container struct {
	ctype      CType // container type
	capacity   int
	mask       uint64 // mask used to assign key to bucket
	seqmask    uint64 // sequence number mask to emulate rollover
	arrcnt     int    // based on ctype
	choices    int
	useCSeqnum bool
	arr        [][]*FifoQ // backing arrays
	seqnum     []uint64   // one per backing array
}

// NewContainer creates a new container of specified CType, with allocated FifoQ array(s)
// as required. Buckets are evenly divided across two arrays for the double array types,
// with key mask adjusted accordingly.
func NewContainer(ctype CType, buckets int, slots int, seqmask uint64) Container {

	var arrcnt int
	var choices int
	var useCSeqnum bool
	switch ctype {
	case BA:
		choices = 1
		arrcnt = 1
	case C2_A1_C:
		choices = 2
		arrcnt = 1
		useCSeqnum = true
	case C2_A1_R:
		choices = 2
		arrcnt = 1
		useCSeqnum = false
	case C2_A2_C:
		choices = 2
		arrcnt = 2
		useCSeqnum = true
	case C2_A2_R:
		choices = 2
		arrcnt = 2
		useCSeqnum = false
	case C4_A4_C:
		choices = 4
		arrcnt = 4
		useCSeqnum = true
	case C4_A4_R:
		choices = 4
		arrcnt = 4
		useCSeqnum = false
	case C8_A8_C:
		choices = 8
		arrcnt = 8
		useCSeqnum = true
	case C8_A8_R:
		choices = 8
		arrcnt = 8
		useCSeqnum = false
	}

	var arrlen = buckets / arrcnt
	var mask = uint64(arrlen - 1)
	c := &container{
		ctype:      ctype,
		capacity:   buckets * slots,
		mask:       mask,
		seqmask:    seqmask,
		arr:        make([][]*FifoQ, arrcnt),
		seqnum:     make([]uint64, arrcnt),
		arrcnt:     arrcnt,
		choices:    choices,
		useCSeqnum: useCSeqnum,
	}
	for i := 0; i < arrcnt; i++ {
		c.arr[i] = make([]*FifoQ, arrlen)
		for j := 0; j < arrlen; j++ {
			c.arr[i][j] = NewFifoQ(slots)
		}
	}
	return c
}

// container.String supports Container.String()
func (p *container) String() string {
	return fmt.Sprintf("type:%s mask:%x seqmask:%x capacity:%d", p.ctype, p.mask, p.seqmask, p.capacity)
}

// container.ArrCnt supports Container.ArrCnt()
func (p *container) ArrCnt() int { return p.arrcnt }

// container.Choices supports Container.Choices()
func (p *container) Choices() int { return p.choices }

// container.UseCSeqnum supports Container.UseCSeqnum
func (p *container) UseCSeqnum() bool { return p.useCSeqnum }

// container.Update supports Container.Update()
func (c *container) Update(p *Params, seqnum uint64, key ...uint64) uint64 {
	if c.ctype == BA {
		idx := int(key[0] & c.mask)
		return c.arr[0][idx].Add(seqnum)
	}

	var choices = c.Choices()
	var arrcnt = c.ArrCnt()
	var useCSeqnum = c.UseCSeqnum()
	var seqnums = make([]uint64, choices)
	var idxs = make([]int, choices)
	for i := 0; i < choices; i++ {
		idxs[i] = int(key[i] & c.mask)
		// arrnum: c2i 0 0 - c2ii 0 1 - c4iv 0 1 2 3
		arrnum := i % arrcnt
		if useCSeqnum {
			seqnums[i] = c.arr[arrnum][idxs[i]].Seqnum()
		} else {
			seqnums[i] = c.arr[arrnum][idxs[i]].Tail()
		}
	}

	pick := PickOldest(p, c.seqmask, seqnum, seqnums)
	idx := idxs[pick]
	arr := c.arr[pick%arrcnt]
	ev := arr[idx].Add(seqnum)
	Trace(p, "evict %x => ", ev)
	arr[idx].DebugPrint(p)
	Trace(p, "------------\n")
	return ev
	//	return arr[idx].Add(seqnum)
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
