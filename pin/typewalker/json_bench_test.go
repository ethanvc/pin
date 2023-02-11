package typewalker

import (
	"encoding/json"
	"testing"
)

func BenchmarkSmallStruct(b *testing.B) {
	type S struct {
		X int
		Y string
	}

	b.Run("pin.ToLogJson", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				v := S{
					X: 3,
					Y: "hello-world",
				}
				ToLogJson(v)
			}
		})
	})

	b.Run("Std.Marshal", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				v := S{
					X: 3,
					Y: "hello-world",
				}
				json.Marshal(v)
			}
		})
	})
}
