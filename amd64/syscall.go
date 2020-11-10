package amd64

func (a *Assembler) Syscall() {
	a.byte(0x0f)
	a.byte(0x05)
}
