package kvrepo

import (
	"fmt"
	"time"
)

// An Kv is a key-value pair.
type Kv struct {
	Key   string
	Value Value
}

// String returns an Kv for a string value.
func String(key, value string) Kv {
	return Kv{key, StringValue(value)}
}

// Int64 returns an Kv for an int64.
func Int64(key string, value int64) Kv {
	return Kv{key, Int64Value(value)}
}

// Int converts an int to an int64 and returns
// an Kv with that value.
func Int(key string, value int) Kv {
	return Int64(key, int64(value))
}

// Uint64 returns an Kv for a uint64.
func Uint64(key string, v uint64) Kv {
	return Kv{key, Uint64Value(v)}
}

// Float64 returns an Kv for a floating-point number.
func Float64(key string, v float64) Kv {
	return Kv{key, Float64Value(v)}
}

// Bool returns an Kv for a bool.
func Bool(key string, v bool) Kv {
	return Kv{key, BoolValue(v)}
}

// Time returns an Kv for a time.Time.
// It discards the monotonic portion.
func Time(key string, v time.Time) Kv {
	return Kv{key, TimeValue(v)}
}

// Duration returns an Kv for a time.Duration.
func Duration(key string, v time.Duration) Kv {
	return Kv{key, DurationValue(v)}
}

// Group returns an Kv for a Group Value.
// The caller must not subsequently mutate the
// argument slice.
//
// Use Group to collect several Kvs under a single
// key on a log line, or as the result of LogValue
// in order to log a single value as multiple Kvs.
func Group(key string, as ...Kv) Kv {
	return Kv{key, GroupValue(as...)}
}

// Any returns an Kv for the supplied value.
// See [Value.AnyValue] for how values are treated.
func Any(key string, value any) Kv {
	return Kv{key, AnyValue(value)}
}

// Equal reports whether a and b have equal keys and values.
func (a Kv) Equal(b Kv) bool {
	return a.Key == b.Key && a.Value.Equal(b.Value)
}

func (a Kv) String() string {
	return fmt.Sprintf("%s=%s", a.Key, a.Value)
}
