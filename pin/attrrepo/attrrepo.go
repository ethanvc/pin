package attrrepo

import "golang.org/x/exp/slog"

const nAttrsInline = 5

type AttrRepo struct {
	// Allocation optimization: an inline array sized to hold
	// the majority of log calls (based on examination of open-source
	// code). It holds the start of the list of Attrs.
	front [nAttrsInline]slog.Attr

	// The number of Attrs in front.
	nFront int

	// The list of Attrs except for those in front.
	// Invariants:
	//   - len(back) > 0 iff nFront == len(front)
	//   - Unused array elements are zero. Used to detect mistakes.
	back []slog.Attr
}

// NumAttrs returns the number of attributes in the Record.
func (r AttrRepo) NumAttrs() int {
	return r.nFront + len(r.back)
}

// Attrs calls f on each Attr in the Record.
// The Attrs are already resolved.
func (r AttrRepo) Attrs(f func(slog.Attr)) {
	for i := 0; i < r.nFront; i++ {
		f(r.front[i])
	}
	for _, a := range r.back {
		f(a)
	}
}

// AddAttrs appends the given Attrs to the Record's list of Attrs.
// It resolves the Attrs before doing so.
func (r *AttrRepo) AddAttrs(attrs ...slog.Attr) {
	resolveAttrs(attrs)
	n := copy(r.front[r.nFront:], attrs)
	r.nFront += n
	// Check if a copy was modified by slicing past the end
	// and seeing if the Attr there is non-zero.
	if cap(r.back) > len(r.back) {
		end := r.back[:len(r.back)+1][len(r.back)]
		if end != (slog.Attr{}) {
			panic("copies of a slog.Record were both modified")
		}
	}
	r.back = append(r.back, attrs[n:]...)
}

// resolveAttrs replaces the values of the given Attrs with their
// resolutions.
func resolveAttrs(as []slog.Attr) {
	for i, a := range as {
		as[i].Value = a.Value.Resolve()
	}
}
