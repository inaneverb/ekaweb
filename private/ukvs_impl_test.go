package ekaweb_private_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/inaneverb/ekaweb/private"
)

func bUkvsInsertStdCtx(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("Insert/StdCtx/%d", size)
	return name, func(b *testing.B) {
		b.ReportAllocs()

		var ctx, cancelFunc = context.WithCancel(context.Background())
		defer cancelFunc()

		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				ctx = context.WithValue(ctx, j, struct{}{})
			}
		}
	}
}

func bUkvsInsertUkvs[M ekaweb_private.UkvsMap](
	name string, gen ekaweb_private.UkvsMapGenerator[M],
	size int) (string, func(b *testing.B)) {

	name = fmt.Sprintf("Insert/Ukvs%s/%d", name, size)

	var keys = make([]any, 0, size)
	for j := 0; j < size; j++ {
		keys = append(keys, strconv.Itoa(j+1))
	}

	var mgr = ekaweb_private.NewUkvsManager(gen, ekaweb_private.RouterOptionCodec{})

	return name, func(b *testing.B) {
		b.ReportAllocs()

		var ctx, cancelFunc = context.WithCancel(context.Background())
		defer cancelFunc()
		ctx = mgr.InjectUkvs(ctx)

		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				ekaweb_private.UkvsInsert(ctx, keys[j], struct{}{})
			}
		}

		mgr.ReturnUkvs(ctx)
	}
}

func bUkvsGetLastStdCtx(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("Get/StdCtx/%d", size)

	var keys = make([]any, 0, size)
	for j := 0; j < size; j++ {
		keys = append(keys, strconv.Itoa(j+1))
	}

	// since context.WithValue() is linked list and stores the last one
	// if we want to get a highest available elem, we need 0'th.
	var key = keys[0]

	return name, func(b *testing.B) {
		b.ReportAllocs()

		var ctx, cancelFunc = context.WithCancel(context.Background())
		defer cancelFunc()

		for j := 0; j < size; j++ {
			ctx = context.WithValue(ctx, keys[j], struct{}{})
		}

		for i := 0; i < b.N; i++ {
			_ = ctx.Value(key)
		}
	}
}

func bUkvsGetLastUkvs[M ekaweb_private.UkvsMap](
	name string, gen ekaweb_private.UkvsMapGenerator[M],
	size int) (string, func(b *testing.B)) {

	name = fmt.Sprintf("Get/Ukvs%s/%d", name, size)

	var keys = make([]any, 0, size)
	for j := 0; j < size; j++ {
		keys = append(keys, strconv.Itoa(j+1))
	}
	var key = keys[len(keys)-1]

	var mgr = ekaweb_private.NewUkvsManager(gen, ekaweb_private.RouterOptionCodec{})

	return name, func(b *testing.B) {
		b.ReportAllocs()

		var ctx, cancelFunc = context.WithCancel(context.Background())
		defer cancelFunc()
		ctx = mgr.InjectUkvs(ctx)

		for j := 0; j < size; j++ {
			ekaweb_private.UkvsInsert(ctx, keys[j], struct{}{})
		}

		for i := 0; i < b.N; i++ {
			ekaweb_private.UkvsLookup(ctx, key)
		}

		mgr.ReturnUkvs(ctx)
	}
}

func bUkvsGetAllStdCtx(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("GetAll/StdCtx/%d", size)

	var keys = make([]any, 0, size)
	for j := 0; j < size; j++ {
		keys = append(keys, strconv.Itoa(j+1))
	}

	return name, func(b *testing.B) {
		b.ReportAllocs()

		var ctx, cancelFunc = context.WithCancel(context.Background())
		defer cancelFunc()

		for j := 0; j < size; j++ {
			ctx = context.WithValue(ctx, keys[j], struct{}{})
		}

		for i := 0; i < b.N; i++ {
			for _, key := range keys {
				_ = ctx.Value(key)
			}
		}
	}
}

func bUkvsGetAllUkvs[M ekaweb_private.UkvsMap](
	name string, gen ekaweb_private.UkvsMapGenerator[M],
	size int) (string, func(b *testing.B)) {

	name = fmt.Sprintf("GetAll/Ukvs%s/%d", name, size)

	var keys = make([]any, 0, size)
	for j := 0; j < size; j++ {
		keys = append(keys, strconv.Itoa(j+1))
	}

	var mgr = ekaweb_private.NewUkvsManager(gen, ekaweb_private.RouterOptionCodec{})

	return name, func(b *testing.B) {
		b.ReportAllocs()

		var ctx, cancelFunc = context.WithCancel(context.Background())
		defer cancelFunc()
		ctx = mgr.InjectUkvs(ctx)

		for j := 0; j < size; j++ {
			ekaweb_private.UkvsInsert(ctx, keys[j], struct{}{})
		}

		for i := 0; i < b.N; i++ {
			for _, key := range keys {
				ekaweb_private.UkvsLookup(ctx, key)
			}
		}

		mgr.ReturnUkvs(ctx)
	}
}

func bUkvsGetHeaderStdCtx(size int) (string, func(b *testing.B)) {
	var name = fmt.Sprintf("GetHeader/StdCtx/%d", size)
	var _, cb = bUkvsGetLastStdCtx(size)
	return name, cb
}

func bUkvsGetHeaderUkvs[M ekaweb_private.UkvsMap](
	name string, gen ekaweb_private.UkvsMapGenerator[M],
	size int) (string, func(b *testing.B)) {

	name = fmt.Sprintf("GetHeader/Ukvs%s/%d", name, size)
	var mgr = ekaweb_private.NewUkvsManager(gen, ekaweb_private.RouterOptionCodec{})

	return name, func(b *testing.B) {
		b.ReportAllocs()

		var ctx, cancelFunc = context.WithCancel(context.Background())
		defer cancelFunc()
		ctx = mgr.InjectUkvs(ctx)
		ekaweb_private.UkvsInsertUserError(ctx, context.Canceled)

		for i := 0; i < b.N; i++ {
			_ = ekaweb_private.UkvsGetUserError(ctx)
		}

		mgr.ReturnUkvs(ctx)
	}
}

func BenchmarkUkvs(b *testing.B) {
	var sizes = []int{1, 2, 4, 8, 12, 16, 24, 32, 64}
	for i, n := 0, len(sizes); i < n; i++ {

		b.Run(bUkvsInsertStdCtx(sizes[i]))
		b.Run(bUkvsGetLastStdCtx(sizes[i]))
		b.Run(bUkvsGetAllStdCtx(sizes[i]))
		b.Run(bUkvsGetHeaderStdCtx(sizes[i]))

		var gen1 = ekaweb_private.NewUkvsMapGeneratorGoMap()
		b.Run(bUkvsInsertUkvs("GoMap", gen1, sizes[i]))
		b.Run(bUkvsGetLastUkvs("GoMap", gen1, sizes[i]))
		b.Run(bUkvsGetAllUkvs("GoMap", gen1, sizes[i]))
		b.Run(bUkvsGetHeaderUkvs("GoMap", gen1, sizes[i]))

		var gen2 = ekaweb_private.NewUkvsMapGeneratorGoMap()
		b.Run(bUkvsInsertUkvs("Slice", gen2, sizes[i]))
		b.Run(bUkvsGetLastUkvs("Slice", gen2, sizes[i]))
		b.Run(bUkvsGetAllUkvs("Slice", gen2, sizes[i]))
		b.Run(bUkvsGetHeaderUkvs("Slice", gen2, sizes[i]))
	}
}
