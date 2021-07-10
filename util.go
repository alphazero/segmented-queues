// Doost
package segque

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"unsafe"
)

// utility functions //////////////////////////////////////////////////////////

/// sort //////////////////////////////////////////////////////////////////////

func ToSortedArrays(m map[float64]float64) ([]float64, []float64) {
	keyset := make([]float64, 0, len(m))
	for k := range m {
		keyset = append(keyset, k)
	}
	sort.Float64s(keyset)
	valset := make([]float64, len(keyset))
	for i, k := range keyset {
		valset[i] = m[k]
	}
	return keyset, valset
}

/// debug and trace emits /////////////////////////////////////////////////////
var w = os.Stdout // REVU less/more deson't work with Stderr
func Trace(p *Params, fmtstr string, v ...interface{}) (int, error) {
	return emit(p.Trace, p, "TRACE "+fmtstr, v...)
}

func Emit(p *Params, fmtstr string, v ...interface{}) (int, error) {
	return emit(p.Verbose, p, fmtstr, v...)
}
func emit(flag bool, p *Params, fmtstr string, v ...interface{}) (int, error) {
	if !flag {
		return 0, nil
	}
	return fmt.Fprintf(w, fmtstr, v...)
}

func ExitOnError(err error, msg string) {
	Exit(1, fmt.Sprintf("%s: error %s\n", msg, err))
}

func Exit(code int, msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(code)
}

/// data files ////////////////////////////////////////////////////////////////

// Creates or truncates existing file. Exits on error.
func CreateDataFile(p *Params) (*os.File, *bufio.Writer) {
	file, err := os.Create(p.Filename)
	if err != nil {
		ExitOnError(err, "CreateDataFile")
	}
	return file, bufio.NewWriter(file)
}
func OpenDataFile(p *Params) (*os.File, *bufio.Reader, int64) {
	file, err := os.Open(p.Filename)
	if err != nil {
		ExitOnError(err, "OpenDataFile")
	}
	finfo, err := file.Stat()
	if err != nil {
		ExitOnError(err, "OpenDataFile")
	}
	return file, bufio.NewReader(file), finfo.Size()
}

func write(op string, w *bufio.Writer, buf []byte) {
	if _, err := w.Write(buf); err != nil {
		ExitOnError(err, op)
	}
}

// Writes the int value - exits on error
func WriteInt(w *bufio.Writer, v int) {
	var buf [8]byte
	buf = *(*[8]byte)(unsafe.Pointer(&v))
	write("WriteInt", w, buf[:])
}

// Writes the float value - exits on error
func WriteUint64(w *bufio.Writer, v uint64) {
	var buf [8]byte
	buf = *(*[8]byte)(unsafe.Pointer(&v))
	write("WriteUint64", w, buf[:])
}

// Writes the float value - exits on error
func WriteFloat64(w *bufio.Writer, v float64) {
	var buf [8]byte
	buf = *(*[8]byte)(unsafe.Pointer(&v))
	write("WriteFloat64", w, buf[:])
}

func read(op string, r *bufio.Reader) ([8]byte, error) {
	var buf [8]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil && err != io.EOF {
		ExitOnError(err, op)
	}
	return buf, err // error can only be nil or EOF
}

// Reads an int value
func ReadInt(r *bufio.Reader) (int, error) {
	p, eof := read("ReadInt", r)
	if eof != nil {
		return 0, eof
	}
	v := *(*int)(unsafe.Pointer(&p))
	return v, nil
}

// Reads an Uint64 value
func ReadUint64(r *bufio.Reader) (uint64, error) {
	p, eof := read("ReadUint64", r)
	if eof != nil {
		return 0, eof
	}
	v := *(*uint64)(unsafe.Pointer(&p))
	return v, nil
}

// Reads a float value
func ReadFloat64(r *bufio.Reader) (float64, error) {
	p, eof := read("ReadFloat64", r)
	if eof != nil {
		return 0, eof
	}
	v := *(*float64)(unsafe.Pointer(&p))
	return v, nil
}
