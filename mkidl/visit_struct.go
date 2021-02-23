package main

import (
	"fmt"

	"github.com/jd3nn1s/gomidl/ast"
)

func (g *generator) generatePublicAndInnerStruct(n *ast.StructNode) {
	publicTypeName := g.toGolangStructType(n.Name, false)
	g.printfln("type %s struct {", publicTypeName)

	simple := true

	for i := 0; i < len(n.Fields)-1; i++ {
		cf := n.Fields[i]

		if g.isInnerOnlyStruct(g.unwrapTypedef(cf.Type)) {
			simple = false
		} else if _, ok := basicTypeMapInner[g.unwrapTypedef(cf.Type)]; ok {
			simple = false
		}

		ct := g.toGolangType(cf.Type, cf.Indirections, false)

		// skip count
		// if strings.HasPrefix(ct, "[]") {
		if ct == "[]string" {
			nf := n.Fields[i+1]
			nt := g.toGolangType(nf.Type, nf.Indirections, false)
			if nt == "uint32" || nt == "int32" {
				continue
			}
		}

		g.printfln("%s %s", cf.Name, ct)
	}

	hasReserved := false
	{
		f := n.Fields[len(n.Fields)-1]

		if f.Name == "Reserved" && f.Type == "void" && f.Indirections == 1 {
			hasReserved = true
		} else {
			// not Reserved
			ct := g.toGolangType(f.Type, f.Indirections, false)
			g.printfln("%s %s", f.Name, ct)
		}
	}
	g.printfln("}")

	innerTypeName := g.toGolangStructType(n.Name, true)

	if simple {

		g.printfln("type %s struct {", innerTypeName)
		g.printfln("%v", publicTypeName)

		if hasReserved {
			g.printfln("Reserved unsafe.Pointer")
		}

		g.printfln("}")

		g.printfln(`func (obj *%s) toGoStruct() *%v {
			return &obj.%v
		}`, innerTypeName, publicTypeName, publicTypeName)

		g.printfln("func (obj *%s) toInnerStruct() *%v {", publicTypeName, innerTypeName)

		if hasReserved {
			g.printfln("return &%v{*obj, nil}", innerTypeName)
		} else {
			g.printfln("return &%v{*obj}", innerTypeName)
		}

		g.printfln("}")

	} else {

		g.printfln("func (obj *%s) toInnerStruct() *%v {", publicTypeName, innerTypeName)
		g.printfln("if obj == nil { return nil}")
		g.printfln("dst := %v{}", innerTypeName)

		for i := 0; i < len(n.Fields)-1; i++ {
			f := n.Fields[i]
			g.generateToInnerObject(fmt.Sprintf("obj.%v", f.Name), fmt.Sprintf("dst.%v", f.Name), f.Type, f.Indirections)
		}

		if !hasReserved {
			f := n.Fields[len(n.Fields)-1]
			g.generateToInnerObject(fmt.Sprintf("obj.%v", f.Name), fmt.Sprintf("dst.%v", f.Name), f.Type, f.Indirections)
		}

		g.printfln("return &dst")
		g.printfln("}")

		g.generateInnerStruct(n)
	}

}

func (g *generator) generateInnerStruct(n *ast.StructNode) {
	publicTypeName := g.toGolangStructType(n.Name, false)
	innerTypeName := g.toGolangStructType(n.Name, true)

	g.printfln("type %s struct {", innerTypeName)
	for _, f := range n.Fields {
		t := g.toGolangType(f.Type, f.Indirections, true)
		g.printfln("%s %s", f.Name, t)

	}
	g.printfln("}")

	if g.isListLike(n.Name) {
		g.printfln("func (obj *%s) toGoStruct() %v {", innerTypeName, publicTypeName)
		g.printfln("var dst %v", publicTypeName)
		g.generateListObjectToGolangSlice("obj", "dst", n.Name)
		g.printfln("return dst")
		g.printfln("}")
	} else if g.isMapLike(n.Name) {
		g.printfln("func (obj *%s) toGoStruct() %v {", innerTypeName, publicTypeName)
		g.printfln("var dst = make(%v)", publicTypeName)
		g.generateMapObjectToGolangMap("obj", "dst", n.Name)
		g.printfln("return dst")
		g.printfln("}")
	} else if !g.isInnerOnlyStruct(n.Name) {
		g.printfln("func (obj *%s) toGoStruct() *%v {", innerTypeName, publicTypeName)
		g.printfln("if obj == nil { return nil}")
		g.printfln("dst := %v{}", publicTypeName)

		for i := 0; i < len(n.Fields)-1; i++ {
			f := n.Fields[i]
			g.generateToGolangObject(fmt.Sprintf("obj.%v", f.Name), fmt.Sprintf("dst.%v", f.Name), f.Type, f.Indirections)
		}

		{
			f := n.Fields[len(n.Fields)-1]
			if !(f.Name == "Reserved" && f.Type == "void" && f.Indirections == 1) {
				f := n.Fields[len(n.Fields)-1]
				g.generateToGolangObject(fmt.Sprintf("obj.%v", f.Name), fmt.Sprintf("dst.%v", f.Name), f.Type, f.Indirections)
			}
		}

		g.printfln("return &dst")
		g.printfln("}")
	}
}

func (g *generator) visitStruct(n *ast.StructNode) {

	if _, ok := basicTypeMap[n.Name]; ok {
		return
	}

	if g.isInnerOnlyStruct(n.Name) {
		g.generateInnerStruct(n)
	} else {
		g.generatePublicAndInnerStruct(n)
	}
}
