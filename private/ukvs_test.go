package ekaweb_private_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/inaneverb/ekaweb/private"
)

type vCtx struct {
	orig context.Context
	kv   vCtxKvPair
}

type vCtxKvPair struct {
	p    *vCtxKvPair
	k, v any
}

func (v vCtx) Deadline() (deadline time.Time, ok bool) {
	return v.orig.Deadline()
}

func (v vCtx) Done() <-chan struct{} {
	return v.orig.Done()
}

func (v vCtx) Err() error {
	return v.orig.Err()
}

func (v vCtx) Value(key any) any {
	var p = &v.kv
	for {
		if p.k == key {
			return p.v
		}
		if p.p == nil {
			return nil
		}
		p = p.p
	}
}

func WithValue(ctx context.Context, k, v any) context.Context {

	var ctx2, ok = ctx.(vCtx)
	if !ok {
		return vCtx{ctx, vCtxKvPair{nil, k, v}}
	}

	var kv = vCtxKvPair{nil, k, v}
	kv.p = ctx2.kv.p
	ctx2.kv.p = &kv

	return ctx2
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func bUkvsInsertStdCtx(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("Insert/StdCtx/%d", size)
	return name, func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var ctx = context.TODO()
			for j := 0; j < size; j++ {
				ctx = context.WithValue(ctx, j, struct{}{})
			}
		}
	}
}

func bUkvsInsertCustomCtx(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("Insert/CustomCtx/%d", size)
	return name, func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var ctx = context.TODO()
			for j := 0; j < size; j++ {
				ctx = WithValue(ctx, j, struct{}{})
			}
		}
	}
}

func bUkvsInsertUkvs(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("Insert/Ukvs/%d", size)
	return name, func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var ctx = ekaweb_private.UkvsInit(context.TODO())
			for j := 0; j < size; j++ {
				ekaweb_private.UkvsInsert(ctx, j, struct{}{})
			}
		}
	}
}

func bUkvsGetStdCtx(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("Get/StdCtx/%d", size)
	var ctx = context.TODO()
	for j := 0; j < size; j++ {
		ctx = context.WithValue(ctx, j, struct{}{})
	}
	return name, func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				_ = ctx.Value(j)
			}
		}
	}
}

func bUkvsGetCustomCtx(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("Get/CustomCtx/%d", size)
	var ctx = context.TODO()
	for j := 0; j < size; j++ {
		ctx = WithValue(ctx, j, struct{}{})
	}
	return name, func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				_ = ctx.Value(j)
			}
		}
	}
}

func bUkvsGetUkvs(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("Get/Ukvs/%d", size)
	var ctx = ekaweb_private.UkvsInit(context.TODO())
	for j := 0; j < size; j++ {
		ekaweb_private.UkvsInsert(ctx, j, struct{}{})
	}
	return name, func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				ekaweb_private.UkvsLookup(ctx, j)
			}
		}
	}
}

func BenchmarkUkvs(b *testing.B) {
	var sizes = []int{1, 2, 4, 8, 12, 16, 24, 32, 64}
	for i, n := 0, len(sizes); i < n; i++ {
		b.Run(bUkvsInsertStdCtx(sizes[i]))
		b.Run(bUkvsInsertCustomCtx(sizes[i]))
		b.Run(bUkvsInsertUkvs(sizes[i]))
		b.Run(bUkvsGetStdCtx(sizes[i]))
		b.Run(bUkvsGetCustomCtx(sizes[i]))
		b.Run(bUkvsGetUkvs(sizes[i]))
	}
}
