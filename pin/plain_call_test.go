package pin

import (
	"context"
	"testing"
)

func BenchmarkPlainCall(b *testing.B) {
	b.Run("CallDirectly", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				req := &TestReq{
					Name: "hello",
				}
				c := &TestController{}
				c.Get(context.Background(), req)
			}
		})
	})

	b.Run("CallWithPlainCall", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				req := &TestReq{
					Name: "hello",
				}
				c := &TestController{}
				CreatePlainCall("Get", c.Get).Call(context.Background(), req)
			}
		})
	})
}
