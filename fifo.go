// Doost!

package segque

type Entry struct {
	V uint64
}

// A queue.
type FifoQ struct {
	seqnum uint64
	cap    int
	idx    int
	buf    []Entry // zero value means never used
}

// NewFifoQ creates a new fifo queue of specified capacity.
// Entries will have sequence numbers in range [0..seqmask]
func NewFifoQ(capacity int) *FifoQ {
	q := FifoQ{
		seqnum: 0,
		cap:    capacity,
		idx:    0,
		buf:    make([]Entry, capacity),
	}
	return &q
}

func (q *FifoQ) Add(v uint64) uint64 {
	q.seqnum = v // container's seqnum is same as last item added
	x := q.buf[q.idx]
	entry := Entry{
		V: v,
	}
	q.buf[q.idx] = entry
	q.idx++
	if q.idx == q.cap {
		q.idx = 0
	}
	return x.V
}

func (q *FifoQ) Seqnum() uint64 {
	return q.seqnum
}

func (q *FifoQ) Tail() uint64 {
	return q.buf[q.idx].V
}

func (q *FifoQ) DebugPrint(p *Params) {
	Trace(p, "%d %x %x\n", q.idx, q.seqnum, q.buf)
}
