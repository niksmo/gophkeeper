package encode

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Encoder struct {
	buf  *bytes.Buffer
	gobE *gob.Encoder
}

func NewEncoder() *Encoder {
	buf := new(bytes.Buffer)
	gobE := gob.NewEncoder(buf)
	return &Encoder{buf, gobE}
}

func (e *Encoder) Encode(src any) ([]byte, error) {
	const op = "encoder.Encode"
	defer e.buf.Reset()
	if err := e.gobE.Encode(src); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	dst := make([]byte, e.buf.Len())
	copy(dst, e.buf.Bytes())
	return dst, nil
}

type Decoder struct {
	buf  *bytes.Buffer
	godD *gob.Decoder
}

func NewDecoder() *Decoder {
	buf := new(bytes.Buffer)
	gobD := gob.NewDecoder(buf)
	return &Decoder{buf, gobD}
}

func (d *Decoder) Decode(dst any, src []byte) error {
	const op = "decoder.Decode"
	d.buf.Write(src)
	defer d.buf.Reset()
	if err := d.godD.Decode(dst); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
