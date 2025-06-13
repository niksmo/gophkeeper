package encode

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/niksmo/gophkeeper/pkg/logger"
)

type Encoder struct {
	l    logger.Logger
	buf  *bytes.Buffer
	gobE *gob.Encoder
}

func NewEncoder(l logger.Logger) *Encoder {
	buf := new(bytes.Buffer)
	gobE := gob.NewEncoder(buf)
	return &Encoder{l, buf, gobE}
}

func (e *Encoder) Encode(src any) ([]byte, error) {
	const op = "encoder.Encode"
	defer e.buf.Reset()
	if err := e.gobE.Encode(src); err != nil {
		e.l.Debug().Str("op", op).Err(err).Msg("failed to encode")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	dst := make([]byte, e.buf.Len())
	copy(dst, e.buf.Bytes())
	return dst, nil
}

type Decoder struct {
	l    logger.Logger
	buf  *bytes.Buffer
	godD *gob.Decoder
}

func NewDecoder(l logger.Logger) *Decoder {
	buf := new(bytes.Buffer)
	gobD := gob.NewDecoder(buf)
	return &Decoder{l, buf, gobD}
}

func (d *Decoder) Decode(dst any, src []byte) error {
	const op = "decoder.Decode"
	d.buf.Write(src)
	defer d.buf.Reset()
	if err := d.godD.Decode(dst); err != nil {
		d.l.Debug().Str("op", op).Err(err).Msg("failed to decode")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
