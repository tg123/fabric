package main

import (
	"sort"
)

func (g *generator) generateStub(callbody func(fn string, paramTypes []string), callbackbody func(fn string, paramTypes []string)) {
	g.importpkg("unsafe")

	sorted := func(m map[int][]string) (keys []int) {
		for k := range m {
			keys = append(keys, k)
		}
		sort.Ints(keys)
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

		_ = p
		// TODO
		// g.printfln("argc int,")
		// g.printfln("this unsafe.Pointer,")

		// for i, a := range p {
		// 	g.printfln("argv%v %v,", i, a)
		// }

		g.printfln(") uintptr {")
		callbackbody(fn, p)
		g.printfln("}")
	}
}

func (g *generator) generateStubDummy() {
	g.headerprintfln(`// +build !windows,!linux,!amd64`)
	g.headerprintfln("")

	returnpanic := func(fn string, paramTypes []string) {
		g.printfln(`panic("not impl")`)
	}

	g.generateStub(returnpanic, returnpanic)
}

func (g *generator) visitStubWindows() {
	g.headerprintfln(`// +build windows,amd64`)
	g.headerprintfln("")

	returnpanic := func(fn string, paramTypes []string) {
		g.printfln(`panic("not impl")`)
	}

	g.generateStub(returnpanic, returnpanic)
}

func (g *generator) visitStubLinux() {
	g.headerprintfln(`// +build linux,amd64`)
	g.headerprintfln("")

	returnpanic := func(fn string, paramTypes []string) {
		g.printfln(`panic("not impl")`)
	}
	g.generateStub(returnpanic, returnpanic)

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
