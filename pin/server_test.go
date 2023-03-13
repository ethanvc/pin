package pin

import (
	"context"
	"testing"
)

func BenchmarkServer_ProcessRequest(b *testing.B) {
	b.Run("ProcessRequest", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				req := &TestReq{}
				c := &TestController{}
				CreatePlainCall("Test", c.Get).Call(context.Background(), req)
			}
		})
	})
}
