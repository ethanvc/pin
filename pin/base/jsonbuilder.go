package base

import "strconv"

type JsonBuilder struct {
	Buf Buffer
}

func (j *JsonBuilder) Reset() {
	j.Buf.Reset()
}

func (j *JsonBuilder) String() string {
	return j.Buf.String()
}

func (j *JsonBuilder) Bytes() []byte {
	return j.Buf.Bytes()
}

func (j *JsonBuilder) OpenArray() *JsonBuilder {
	j.Buf.WriteByte('[')
	return j
}

func (j *JsonBuilder) CloseArray() *JsonBuilder {
	j.removeComma()
	j.Buf.WriteByte(']')
	j.writeComma()
	return j
}

func (j *JsonBuilder) OpenObject() *JsonBuilder {
	j.Buf.Write([]byte("{"))
	return j
}

func (j *JsonBuilder) CloseObject() *JsonBuilder {
	j.removeComma()
	j.Buf.Write([]byte("}"))
	j.writeComma()
	return j
}

func (j *JsonBuilder) WriteKey(key string) *JsonBuilder {
	j.Buf.WriteByte('"')
	j.Buf.WriteString(key)
	j.Buf.WriteByte('"')
	j.Buf.WriteByte(':')
	return j
}

func (j *JsonBuilder) WriteValueBool(v bool) *JsonBuilder {
	return j
}

func (j *JsonBuilder) WriteValueNull() *JsonBuilder {
	j.Buf.Write([]byte("null"))
	j.writeComma()
	return j
}

func (j *JsonBuilder) WriteValueInt64(v int64) *JsonBuilder {
	s := strconv.FormatInt(v, 10)
	j.Buf.WriteString(s)
	j.writeComma()
	return j
}

func (j *JsonBuilder) WriteValueUint64(v uint64) *JsonBuilder {
	s := strconv.FormatUint(v, 10)
	j.Buf.WriteString(s)
	j.writeComma()
	return j
}

func (j *JsonBuilder) WriteValueFloat(v float64, bitSize int) *JsonBuilder {
	s := strconv.FormatFloat(v, 'f', -1, bitSize)
	j.Buf.WriteString(s)
	j.writeComma()
	return j
}

func (j *JsonBuilder) WriteValueString(s string) *JsonBuilder {
	j.Buf.WriteByte('"')
	j.Buf.WriteString(s)
	j.Buf.WriteByte('"')
	j.writeComma()
	return j
}

func (j *JsonBuilder) Finish() {
	j.removeComma()
}

func (j *JsonBuilder) writeComma() {
	j.Buf.WriteByte(',')
}

func (j *JsonBuilder) removeComma() {
	j.Buf.TrimRight(",")
}

func (j *JsonBuilder) writeString(s string) {
	j.Buf.WriteString(s)
}
