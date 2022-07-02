// go test -v types.go types_test.go

package send_data

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType     uint8  = iota + 1 // 1: 1 byte header
	StringType                       // 2: 2
	MaxPayloadSize uint32 = 10 << 20 // 10485760: 10 MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

////////////////////////////////////////////////

type Binary []byte

func (m Binary) Bytes() []byte  { return m }
func (m Binary) String() string { return string(m) }

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	// Write: 1 byte header
	err := binary.Write(w, binary.BigEndian, BinaryType)
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	// Write: Payload Size
	err = binary.Write(w, binary.BigEndian, uint32(len(m)))
	if err != nil {
		return n, err
	}
	n += 4

	// Write: Payload
	o, err := w.Write(m) // o == len(m)

	// 1 + 4 + payload size
	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	// Read: 1 byte header
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ) // 1 byte type
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != BinaryType {
		return n, errors.New("invalid Binary")
	}

	// Read: Payload Size
	var size uint32
	err = binary.Read(r, binary.BigEndian, &size) // 4 bytes size
	if err != nil {
		return n, err
	}
	n += 4
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	// Read: Payload
	*m = make([]byte, size)
	o, err := r.Read(*m) // payload

	// 1 + 4 + payload size
	return n + int64(o), err
}

////////////////////////////////////////////////

type String string

func (m String) Bytes() []byte  { return []byte(m) }
func (m String) String() string { return string(m) }

func (m String) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, StringType) // 1 byte type
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m)))
	if err != nil {
		return n, err
	}
	n += 4

	o, err := w.Write([]byte(m)) // payload

	return n + int64(o), err
}

func (m *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ) // 1 byte type
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != StringType {
		return n, errors.New("invalid String")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size) // 4 bytes size
	if err != nil {
		return n, err
	}

	buf := make([]byte, size)
	o, err := r.Read(buf) // payload
	if err != nil {
		return n, err
	}

	*m = String(buf)

	return n + int64(o), nil
}

////////////////////////////////////////////////

func decode(r io.Reader) (Payload, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	var payload Payload
	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StringType:
		payload = new(String)
	default:
		return nil, errors.New("unkown type")
	}

	_, err = payload.ReadFrom(io.MultiReader(bytes.NewReader([]byte{typ}), r))
	if err != nil {
		return nil, err
	}

	return payload, nil

}

/*

=== RUN   TestPayloads
    types_test.go:58: [*send_data.Binary] "Clear is better than clever."
    types_test.go:58: [*send_data.String] "Errors are values."
    types_test.go:58: [*send_data.Binary] "Don't panic."
--- PASS: TestPayloads (0.00s)
=== RUN   TestMaxPayloadSize
--- PASS: TestMaxPayloadSize (0.00s)
PASS
ok      command-line-arguments  0.563s

*/
