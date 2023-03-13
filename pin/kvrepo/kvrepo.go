package kvrepo

const nKvsInline = 5

type KvRepo struct {
	// Allocation optimization: an inline array sized to hold
	// the majority of log calls (based on examination of open-source
	// code). It holds the start of the list of Kvs.
	front [nKvsInline]Kv

	// The number of Kvs in front.
	nFront int

	// The list of Kvs except for those in front.
	// Invariants:
	//   - len(back) > 0 iff nFront == len(front)
	//   - Unused array elements are zero. Used to detect mistakes.
	back []Kv
}

// NumKvs returns the number of attributes in the Record.
func (r KvRepo) NumKvs() int {
	return r.nFront + len(r.back)
}

// Kvs calls f on each Attr in the Record.
// The Kvs are already resolved.
func (r KvRepo) Kvs(f func(Kv)) {
	for i := 0; i < r.nFront; i++ {
		f(r.front[i])
	}
	for _, a := range r.back {
		f(a)
	}
}

// AddKvs appends the given Kvs to the Record's list of Kvs.
// It resolves the Kvs before doing so.
func (r *KvRepo) AddKvs(attrs ...Kv) {
	resolveKvs(attrs)
	n := copy(r.front[r.nFront:], attrs)
	r.nFront += n
	// Check if a copy was modified by slicing past the end
	// and seeing if the Attr there is non-zero.
	if cap(r.back) > len(r.back) {
		end := r.back[:len(r.back)+1][len(r.back)]
		if end != (Kv{}) {
			panic("copies of a slog.Record were both modified")
		}
	}
	r.back = append(r.back, attrs[n:]...)
}

// resolveKvs replaces the values of the given Kvs with their
// resolutions.
func resolveKvs(as []Kv) {
	for i, a := range as {
		as[i].Value = a.Value.Resolve()
	}
}
