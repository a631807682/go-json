package json

import (
	"errors"
	"unsafe"
)

type stringDecoder struct {
}

func newStringDecoder() *stringDecoder {
	return &stringDecoder{}
}

func (d *stringDecoder) decode(ctx *context, p uintptr) error {
	bytes, err := d.decodeByte(ctx)
	if err != nil {
		return err
	}
	*(*string)(unsafe.Pointer(p)) = *(*string)(unsafe.Pointer(&bytes))
	return nil
}

func (d *stringDecoder) decodeByte(ctx *context) ([]byte, error) {
	buf := ctx.buf
	buflen := ctx.buflen
	cursor := ctx.cursor
	for ; cursor < buflen; cursor++ {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			continue
		case '"':
			cursor++
			start := cursor
			for ; cursor < buflen; cursor++ {
				tk := buf[cursor]
				if tk == '\\' {
					cursor++
					continue
				}
				if tk == '"' {
					literal := buf[start:cursor]
					cursor++
					ctx.cursor = cursor
					return literal, nil
				}
			}
			return nil, errors.New("unexpected error string")
		}
	}
	return nil, errors.New("unexpected error key delimiter")
}
