package main

import (
	"fmt"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/jd3nn1s/gomidl/scanner"
)

func extractSizeof(n *ast.StructNode) (fieldSize map[string]string, fieldIsSizeOf map[string]string) {
	fieldSize = make(map[string]string)
	fieldIsSizeOf = make(map[string]string)
	for _, f := range n.Fields {
		for _, a := range f.Attributes {
			if a.Type == scanner.SIZE_IS {
				fieldSize[f.Name] = a.Val
				fieldIsSizeOf[a.Val] = f.Name
			}
		}
	}

	return
}

func (g *generator) generatePublicStructFields(n *ast.StructNode) (simple bool) {
	simple = true
	fieldSize, fieldIsSizeOf := extractSizeof(n)

	for i := 0; i < len(n.Fields)-1; i++ {
		cf := n.Fields[i]

		if g.isInnerOnlyStruct(g.unwrapTypedef(cf.Type)) {
			simple = false
		} else if _, ok := basicTypeMapInner[g.unwrapTypedef(cf.Type)]; ok {
			simple = false
		}

		// skip count
		if _, ok := fieldIsSizeOf[cf.Name]; ok {
			continue
		}

		var ct string
		if _, ok := fieldSize[cf.Name]; ok {
			ct = "[]" + g.toGolangType(cf.Type, cf.Indirections-1, false)
		} else {
			ct = g.toGolangType(cf.Type, cf.Indirections, false)
		}

		g.printfln("%s %s", cf.Name, ct)
	}

	{
		if !isResversed(n, len(n.Fields)-1) {
			f := n.Fields[len(n.Fields)-1]
			ct := g.toGolangType(f.Type, f.Indirections, false)
			g.printfln("%s %s", f.Name, ct)

		}
	}

	return
}
func (g *generator) generateGolangToInner(src, dst string, n *ast.StructNode) {
	fieldSize, fieldIsSizeOf := extractSizeof(n)

	for i := 0; i < len(n.Fields)-1; i++ {
		f := n.Fields[i]
		sizeof := fieldIsSizeOf[f.Name]
		if sizeof != "" {
			continue
		}

		sizeis := fieldSize[f.Name]

		if sizeis != "" {
			g.generateSliceToPointer(fmt.Sprintf("%v.%v", src, f.Name), fmt.Sprintf("%v.%v", dst, sizeis), fmt.Sprintf("%v.%v", dst, f.Name), f.Type, f.Indirections-1)
			continue
		}

		g.generateToInnerObject(fmt.Sprintf("%v.%v", src, f.Name), fmt.Sprintf("%v.%v", dst, f.Name), f.Type, f.Indirections)
	}

	if !isResversed(n, len(n.Fields)-1) {
		f := n.Fields[len(n.Fields)-1]
		g.generateToInnerObject(fmt.Sprintf("%v.%v", src, f.Name), fmt.Sprintf("%v.%v", dst, f.Name), f.Type, f.Indirections)
	}
}

func (g *generator) generatePublicAndInnerStruct(n *ast.StructNode) {
	publicTypeName := g.toGolangStructType(n.Name, false)
	g.printfln("type %s struct {", publicTypeName)

	simple := g.generatePublicStructFields(n)

	for _, c := range g.ctx.definedStructEx[n.Name] {
		simple = false
		g.generatePublicStructFields(g.ctx.definedStruct[c])
	}

	g.printfln("}")

	hasReserved := isResversed(n, len(n.Fields)-1)

	innerTypeName := g.toGolangStructType(n.Name, true)

	if simple {

		g.printfln("type %s struct {", innerTypeName)
		g.printfln("%v", publicTypeName)

		if hasReserved {
			g.printfln("Reserved unsafe.Pointer")
		}

		g.printfln("}")

		g.printfln(`func (obj *%s) toGoStruct() *%v {
			if obj == nil { return nil }
			return &obj.%v
		}`, innerTypeName, publicTypeName, publicTypeName)

		g.printfln("func (obj *%s) toInnerStruct() *%v {", publicTypeName, innerTypeName)

		g.printfln(`if obj == nil { return nil }`)

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
		g.generateGolangToInner("obj", "dst", n)

		if len(g.ctx.definedStructEx[n.Name]) > 0 {
			g.printfln("ex0 := dst")
		}

		g.importpkg("unsafe")
		for i, c := range g.ctx.definedStructEx[n.Name] {
			g.printfln("ex%v := &%v{}", i+1, g.toGolangStructType(c, true))
			g.generateGolangToInner("obj", fmt.Sprintf("ex%v", i+1), g.ctx.definedStruct[c])

			g.printfln("ex%v.Reserved = unsafe.Pointer(ex%v)", i, i+1)
		}

		g.printfln("return &dst")
		g.printfln("}")

		g.generateInnerStructAndConverter(n)
	}

}

func (g *generator) generateInnerStruct(n *ast.StructNode) {
	innerTypeName := g.toGolangStructType(n.Name, true)

	g.printfln("type %s struct {", innerTypeName)
	for _, f := range n.Fields {
		t := g.toGolangType(f.Type, f.Indirections, true)
		g.printfln("%s %s", f.Name, t)
	}
	g.printfln("}")
}

func isResversed(n *ast.StructNode, i int) bool {
	if i != len(n.Fields)-1 {
		return false
	}
	f := n.Fields[i]
	return (f.Name == "Reserved" && f.Type == "void" && f.Indirections == 1)
}

func (g *generator) generateInnerToGo(src, dst string, n *ast.StructNode) {
	fieldSize, fieldIsSizeOf := extractSizeof(n)
	for i := 0; i < len(n.Fields)-1; i++ {
		f := n.Fields[i]
		if _, ok := fieldIsSizeOf[f.Name]; ok {
			continue
		}

		sizeis := fieldSize[f.Name]
		if sizeis != "" {
			g.generatePointerToGolangSlice(fmt.Sprintf("%v.%v", dst, f.Name), fmt.Sprintf("%v.%v", src, sizeis), fmt.Sprintf("%v.%v", src, f.Name), f.Type, f.Indirections-1)
			continue
		}

		g.generateToGolangObject(fmt.Sprintf("%v.%v", src, f.Name), fmt.Sprintf("%v.%v", dst, f.Name), f.Type, f.Indirections)
	}

	if !isResversed(n, len(n.Fields)-1) {
		f := n.Fields[len(n.Fields)-1]
		g.generateToGolangObject(fmt.Sprintf("%v.%v", src, f.Name), fmt.Sprintf("%v.%v", dst, f.Name), f.Type, f.Indirections)
	}
}

func (g *generator) generateInnerStructAndConverter(n *ast.StructNode) {
	publicTypeName := g.toGolangStructType(n.Name, false)
	innerTypeName := g.toGolangStructType(n.Name, true)

	g.generateInnerStruct(n)

	for _, ex := range g.ctx.definedStructEx[n.Name] {
		g.generateInnerStruct(g.ctx.definedStruct[ex])
	}

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

		g.printfln("func (obj *%s) toGoStruct() (dst *%v) {", innerTypeName, publicTypeName)
		g.printfln("if obj == nil { return}")
		g.printfln("dst = &%v{}", publicTypeName)
		g.generateInnerToGo("obj", "dst", n)
		if len(g.ctx.definedStructEx[n.Name]) > 0 {
			g.printfln("ex0 := obj")
		}

		for i, c := range g.ctx.definedStructEx[n.Name] {
			g.printfln("ex%v := (*%v)(ex%v.Reserved)", i+1, g.toGolangStructType(c, true), i)
			g.printfln("if ex%v == nil { return }", i+1)
			g.generateInnerToGo(fmt.Sprintf("ex%v", i+1), "dst", g.ctx.definedStruct[c])
		}

		g.printfln("return")
		g.printfln("}")
	}
}

func (g *generator) visitStruct(n *ast.StructNode) {

	if _, ok := basicTypeMap[n.Name]; ok {
		return
	}

	if g.ctx.definedStructExParent[n.Name] != "" {
		return
	}

	if g.isInnerOnlyStruct(n.Name) {
		g.generateInnerStructAndConverter(n)
	} else {
		g.generatePublicAndInnerStruct(n)
	}
}
