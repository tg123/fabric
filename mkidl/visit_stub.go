package main

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

func (g *generator) generateStub(callbody func(fn string, paramTypes []string), callbackbody func(fn string, paramTypes []string)) {
	g.importpkg("unsafe")

	sorted := func(m map[uint32][]string) (keys []uint32) {
		for k := range m {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		return
	}

	for _, f := range sorted(g.ctx.stubBuilder.calls) {

		fn := g.ctx.stubBuilder.CallStubName(f)
		p := g.ctx.stubBuilder.calls[f]

		g.printfln("func %v(", fn)
		g.printfln("addr uintptr,")
		g.printfln("argc int,")
		g.printfln("this unsafe.Pointer,")

		for i, a := range p {
			ap := strings.Split(a, ".")
			if len(ap) == 2 {
				g.importpkg(strings.Trim(ap[0], "*"))
			}
			g.printfln("argv%v %v,", i, a)
		}

		g.printfln(") (uintptr, error) {")
		callbody(fn, p)
		g.printfln("}")
	}

	for _, f := range sorted(g.ctx.stubBuilder.callbacks) {

		fn := g.ctx.stubBuilder.CallbackStubName(f)
		p := g.ctx.stubBuilder.calls[f]

		g.printfln("func %v(", fn)
		g.printfln("cb interface{},")
		g.printfln(") uintptr {")
		callbackbody(fn, p)
		g.printfln("}")
	}
}

func (g *generator) generateStubDummy() {
	g.headerprintfln(`// +build !windows,!linux,!amd64,!cgo`)
	g.headerprintfln("")

	returnpanic := func(fn string, paramTypes []string) {
		g.printfln(`panic("not impl")`)
	}

	g.generateStub(returnpanic, returnpanic)
}

func (g *generator) visitStubWindows() {
	g.headerprintfln(`// +build windows,amd64`)
	g.headerprintfln("")

	g.generateStub(func(fn string, paramTypes []string) {
		numSyscallParams := len(paramTypes)
		numSyscallRequired := int(math.Ceil(float64(numSyscallParams+1)/3)) * 3

		syscallFunc := "Syscall"
		if numSyscallRequired > 18 {
			panic("more params than supported by syscall")
		} else if numSyscallRequired > 3 {
			syscallFunc += strconv.Itoa(numSyscallRequired)
		}

		g.importpkg("unsafe")
		g.importpkg("syscall")

		g.printfln("hr, _, err := syscall.%s(", syscallFunc)
		g.printfln("addr,")
		g.printfln("uintptr(argc + 1),")
		g.printfln("uintptr(this),")

		for i := 0; i < numSyscallRequired-1; i++ {
			if i < numSyscallParams {

				if strings.HasPrefix(paramTypes[i], "*") {
					g.printfln("uintptr(unsafe.Pointer(argv%v)),", i)
					continue
				}

				switch paramTypes[i] {
				case "bool":
					g.printfln("boolToUintptr(argv%v),", i)
				case "ole.GUID":
					g.printfln("uintptr(unsafe.Pointer(argv%v)),", i)
				default:
					g.printfln("uintptr(argv%v),", i)
				}

			} else {
				g.printfln("0,")
			}
		}

		g.printfln(")")
		g.printfln("return hr, err")

	}, func(fn string, paramTypes []string) {
		g.importpkg("syscall")

		g.printfln("return syscall.NewCallback(cb)")

	})
}

func (g *generator) visitStubLinux() {
	g.headerprintfln(`// +build linux,amd64,cgo`)
	g.headerprintfln("")

	returnpanic := func(fn string, paramTypes []string) {
		g.printfln(`panic("not impl")`)
	}
	g.generateStub(returnpanic, func(fn string, paramTypes []string) {
		g.printfln(`panic("not impl")`)
	})

}

func (g *generator) visitStub(s string) {
	switch s {
	case "dummy":
		g.generateStubDummy()
	case "windows":
		g.visitStubWindows()
	case "linux":
		g.visitStubLinux()
	}
}
