package base

func In[T comparable](v T, others ...T) bool {
	for _, other := range others {
		if v == other {
			return true
		}
	}
	return false
}

func NotIn[T comparable](v T, others ...T) bool {
	return !In(v, others...)
}
