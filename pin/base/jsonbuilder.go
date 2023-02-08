package base

type JsonBuilder struct {
	Buf      Buffer
	commaSep []byte
}

func (j *JsonBuilder) String() string {
	return j.Buf.String()
}

func (j *JsonBuilder) Bytes() []byte {
	return j.Buf.Bytes()
}

func (j *JsonBuilder) OpenArray() *JsonBuilder {
	return j
}

func (j *JsonBuilder) CloseArray() *JsonBuilder {
	return j
}

func (j *JsonBuilder) OpenObject() *JsonBuilder {
	return j
}

func (j *JsonBuilder) CloseObject() *JsonBuilder {
	return j
}

func (j *JsonBuilder) WriteKey(key string) *JsonBuilder {
	return j
}

func (j *JsonBuilder) WriteValueBool(v bool) *JsonBuilder {
	return j
}

func (j *JsonBuilder) WriteValueNull() *JsonBuilder {
	j.beforeWriteValue()
	j.Buf.Write([]byte("null"))
	return j
}

func (j *JsonBuilder) beforeWriteValue() *JsonBuilder {
	j.Buf.Write(j.commaSep)
	return j
}
