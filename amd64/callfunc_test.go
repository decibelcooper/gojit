package amd64

import (
	"testing"

	"github.com/nelhage/gojit"
)

func TestCallFunc(t *testing.T) {
	asm := newAsm(t)
	defer gojit.Release(asm.Buf)

	called := false

	asm.CallFunc(func() { called = true })
	asm.Ret()

	var f func()
	asm.BuildTo(&f)
	f()

	if !called {
		t.Error("CallFunc did not call the function")
	}
}

func TestRecursion(t *testing.T) {
	asm := newAsm(t)
	defer gojit.Release(asm.Buf)

	var jitf func(i int)
	gof := func(i int) {
		if i > 0 {
			jitf(i - 1)
		}
	}

	asm.Mov(Indirect{Rdi, 0, 64}, Rax)
	asm.Push(Rax)
	asm.CallFunc(gof)
	asm.Pop(Rax)
	asm.Ret()

	asm.BuildTo(&jitf)

	jitf(16)
}

func BenchmarkGoCall(b *testing.B) {
	asm, _ := NewGoABI(gojit.PageSize)
	defer asm.Release()

	f := func() {}
	asm.CallFunc(f)
	asm.Ret()

	var jit func()
	asm.BuildTo(&jit)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jit()
	}
}
