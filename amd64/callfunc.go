package amd64

import (
	"reflect"
	"unsafe"
)

func (a *Assembler) CallFunc(f interface{}) {
	switch a.ABI {
	case GoABI:
		a.CallFuncGo(f)
	default:
		panic("bad ABI")
	}
}

// CallFuncGo assembles a call directly to the go function 'f'. No stack
// swtitching or other setup is performed.
func (a *Assembler) CallFuncGo(f interface{}) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		panic("CallFunc: Can't call non-func")
	}
	ival := *(*struct {
		typ uintptr
		fun uintptr
	})(unsafe.Pointer(&f))
	a.MovAbs(uint64(ival.fun), Rdx)
	a.Call(Indirect{Rdx, 0, 64})
}
