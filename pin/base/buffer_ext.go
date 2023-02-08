package base

import "bytes"

func (b *Buffer) TrimRight(cutset string) *Buffer {
	b.buf = bytes.TrimRight(b.buf, cutset)
	return b
}
