// Doost!

package segque

import (
	"fmt"
	"github.com/minio/blake2b-simd"
	"hash/maphash"
	"time"
	"unsafe"
)

type HashFunc func(b []byte, param interface{}) uint64

type FuncParam func(int) interface{}

type HType int

const (
	_ HType = iota
	Blake2b
	GomapHash
)

func (t HType) String() string {
	return hashtypes[t]
}

type htypeName struct {
	htype HType
	name  string
}

var hashtypes = map[HType]string{
	Blake2b:   "Blake2b",
	GomapHash: "GomapHash",
}

func NewHashFunc(htype HType) (HashFunc, FuncParam) {
	var hfunc HashFunc
	var pfunc FuncParam

	switch htype {
	case Blake2b:
		hfunc = B2bHash
		pfunc = B2bParam
	case GomapHash:
		hfunc = GmHash
		pfunc = GmParam
	}
	return hfunc, pfunc
}

/// rand value generators ////////////////////////////////////////

var nonce uint64

// Randu64 returns a random uint64 using the provided hashfunction,
// the value v, and optional params for the hashfunc.
func Randu64(hf HashFunc, v int, param interface{}) uint64 {
	nonce++
	t := time.Now().UnixNano()
	s := fmt.Sprintf("%016 %016 %016x", t, nonce, v)
	return hf([]byte(s), param)
}

func B2bParam(x int) interface{} { return x }
func B2bHash(b []byte, param interface{}) uint64 {
	var idx = param.(int)
	h := blake2b.Sum256(b)
	return *(*uint64)(unsafe.Pointer(&h[idx]))
}

func GmParam(x int) interface{} { return maphash.MakeSeed() }
func GmHash(b []byte, param interface{}) uint64 {
	var seed = param.(maphash.Seed)
	var h maphash.Hash
	h.SetSeed(seed)
	h.Write(b)
	return h.Sum64()
}
