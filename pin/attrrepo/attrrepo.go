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

func (r AttrRepo) Attrs(f func(slog.Attr)) {
	for i := 0; i < r.nFront; i++ {
		f(r.front[i])
	}
	for _, a := range r.back {
		f(a)
	}
}
