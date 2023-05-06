package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

const (
	DatagramSize = 516              // maximum size of a datagram
	BlockSize    = DatagramSize - 4 // 4 bytes for header
)

type OpCode uint16

const (
	OpRRQ OpCode = iota + 1
	_            // WRQ
	OpData
	OpAck
	OpError
)

type ErrCode uint16

const (
	ErrUnknown ErrCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)

type ReadReq struct {
	Filename string
	Mode     string
}

func (q ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if q.Mode != "" {
		mode = q.Mode
	}

	// OP Code + Filename + 0 + Mode + 0
	cap := 2 + 2 + len(q.Filename) + 1 + len(mode) + 1

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpRRQ)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(q.Filename)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil

}

func (q *ReadReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpRRQ {
		return errors.New("invalid opcode")
	}

	q.Filename, err = r.ReadString(0) // read filename
	if err != nil {
		return errors.New("invalid RRQ")
	}

	q.Filename = strings.TrimRight(q.Filename, "\x00") // remove trailing 0
	if len(q.Filename) == 0 {
		return errors.New("invalid RRQ")
	}

	q.Mode, err = r.ReadString(0) // read mode
	if err != nil {
		return errors.New("invalid RRQ")
	}

	q.Mode = strings.TrimRight(q.Mode, "\x00") // remove trailing 0
	if len(q.Mode) == 0 {
		return errors.New("invalid RRQ")
	}

	actual := strings.ToLower(q.Mode) // convert to lower case
	if actual != "octet" {
		return errors.New("only binary transfers are supported")
	}

	return nil
}

type Data struct {
	Block   uint16
	Payload io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	b.Grow(DatagramSize)

	d.Block++ // increment block number

	err := binary.Write(b, binary.BigEndian, OpData)
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, d.Block)
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(b, d.Payload, BlockSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *Data) UnmarshalBinary(p []byte) error {

	if l := len(p); l < 4 || l > DatagramSize {
		return errors.New("invalid Data")
	}

	var opcode OpCode

	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &opcode)
	if err != nil || opcode != OpData {
		return errors.New("invalid Data")
	}

	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block)
	if err != nil {
		return errors.New("invalid Data")
	}

	d.Payload = bytes.NewReader(p[4:])

	return nil
}

type Ack uint16

func (a Ack) MarshalBinary() ([]byte, error) {
	cap := 2 + 2 // OP Code + Block

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpAck) // write OP Code
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, a) // write block number
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (a *Ack) UnmarshalBinary(p []byte) error {
	var code OpCode

	r := bytes.NewReader(p)

	err := binary.Read(r, binary.BigEndian, &code) // read OP Code
	if err != nil {
		return err
	}

	if code != OpAck {
		return errors.New("invalid ACK")
	}

	return binary.Read(r, binary.BigEndian, a) // read block number
}

type Error struct {
	Error   ErrCode
	Message string
}

func (e Error) MarshalBinary() ([]byte, error) {
	// OP Code + Error Code + Message + 0
	cap := 2 + 2 + len(e.Message) + 1

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpError) // write OP Code
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, e.Error) // write Error Code
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(e.Message) // write message
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // write trailing 0
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (e *Error) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code) // read OP Code
	if err != nil {
		return err
	}

	if code != OpError {
		return errors.New("invalid Error")
	}

	err = binary.Read(r, binary.BigEndian, &e.Error) // read Error Code
	if err != nil {
		return err
	}

	e.Message, err = r.ReadString(0)                 // read message
	e.Message = strings.TrimRight(e.Message, "\x00") // remove trailing 0

	return err
}
