package base

type JsonBuilder struct {
	buf Buffer
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
	return j
}
