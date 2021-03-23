package main

import (
	"fmt"
	"strings"

	"github.com/pinzolo/casee"
)

var basicTypeMap = map[string]string{
	"BOOLEAN":   "bool",
	"BOOL":      "bool",
	"DOUBLE":    "float64",
	"double":    "float64",
	"LONG":      "int32",
	"DWORD":     "uint32",
	"ULONG":     "uint32",
	"ULONGLONG": "uint64",
	"LONGLONG":  "int64",
	"BYTE":      "byte",
	"byte":      "byte",
	"int":       "int32",
	"LPWSTR":    "string",
	"LPCWSTR":   "string",

	"HRESULT":  "FabricErrorCode",
	"GUID":     "uuid.UUID",
	"FILETIME": "time.Time",
	"void":     "interface{}",
}

var basicTypeMapInner = map[string]string{
	"LPWSTR":   "*uint16",
	"LPCWSTR":  "*uint16",
	"FILETIME": "filetime",
	"void":     "unsafe.Pointer",
}

func (g *generator) unwrapTypedef(idltype string) string {
	if t, ok := g.ctx.definedTypedef[idltype]; ok {
		return t.Type
	}

	return idltype
}

var innerOnlyWhitelist = map[string]bool{
	"FABRIC_NAMED_PARTITION_SCHEME_DESCRIPTION": true,
	"FABRIC_NAMED_REPARTITION_DESCRIPTION":      true,
}

func (g *generator) isListLike(idltype string) bool {
	if strings.HasSuffix(idltype, "_LIST") {
		return true
	}

	return false
}

func (g *generator) isMapLike(idltype string) bool {
	return strings.HasSuffix(idltype, "_MAP")
}

func (g *generator) isInnerOnlyStruct(idltype string) bool {
	idltype = g.unwrapTypedef(idltype)

	if g.isListLike(idltype) {
		return true
	}

	if g.isMapLike(idltype) {
		return true
	}

	return innerOnlyWhitelist[idltype]
}

func (g *generator) toGolangType(idltype string, indirections int, inner bool) string {
	// TODO bug indirections not support by unwrapTypedef
	// idltype = g.unwrapTypedef(idltype)
	if t, ok := g.ctx.definedTypedef[idltype]; ok {
		idltype = t.Type
		indirections += t.Indirection
	}

	gotype := basicTypeMap[idltype]

	if gotype != "" {
		if inner {
			gotype = basicTypeMapInner[idltype]
		}

		if gotype == "" {
			gotype = basicTypeMap[idltype]
		}

		p := strings.Split(gotype, ".")

		if len(p) == 2 {
			g.importpkg(p[0])
		}

		// hacks remove later
		switch gotype {
		case "unsafe.Pointer":
			fallthrough
		case "interface{}":
			indirections = 0
		case "string":
			if indirections == 1 {
				gotype = "[]string"
				// } else {
				// panic(fmt.Sprintf("*string %v %v", indirections, idltype))
			}
		}

	} else if _, ok := g.ctx.definedStruct[idltype]; ok {
		// a struct
		gotype = g.toGolangStructType(idltype, inner)
	} else if _, ok := g.ctx.definedEnum[idltype]; ok {
		gotype = goEnumName(idltype)
	} else if _, ok := g.ctx.definedInterface[idltype]; ok {
		gotype = g.goInterfaceName(idltype)
	}

	if gotype == "" {
		panic(fmt.Sprintf("%v unmapped type", idltype))
	}

	if strings.HasPrefix(gotype, "[]") || strings.HasPrefix(gotype, "map[") {
		indirections -= 1
	}

	return fmt.Sprintf("%v%v", strings.Repeat("*", indirections), gotype)
}

func (g *generator) toGolangStructType(name string, inner bool) string {
	name = g.unwrapTypedef(name)

	if inner {
		return fmt.Sprintf("inner%v", casee.ToPascalCase(name))
	}

	structDef := g.ctx.definedStruct[name]

	// list to slice
	if strings.HasSuffix(name, "_LIST") {
		f := structDef.Fields[1] // this is the items def

		// assert
		for {
			if len(structDef.Fields) == 2 && f.Indirections == 1 {
				break
			}

			// this one has an extra reversed
			if name == "FABRIC_APPLICATION_LOAD_METRIC_INFORMATION_LIST" {
				break
			}

			// this one has a total and a reversed
			if strings.HasSuffix(name, "CHUNK_LIST") && len(structDef.Fields) == 4 {
				break
			}

			panic(fmt.Sprintf("%v bad list", name))
		}

		return fmt.Sprintf("[]%v", g.toGolangType(f.Type, 0, false))

	} else if strings.HasSuffix(name, "_MAP") {
		f := structDef.Fields[1] // this is the items

		if len(structDef.Fields) != 2 {
			panic(fmt.Sprintf("%v bad map", name))
		}

		if f.Name != "Items" {
			panic(fmt.Sprintf("%v bad map item", name))
		}

		kvt := g.ctx.definedStruct[g.unwrapTypedef(f.Type)]

		// if len(kvt.Fields) == 2 {
		// 	panic(fmt.Sprintf("%v bad map item type", name))
		// }

		return fmt.Sprintf("map[%v]%v", g.toGolangType(kvt.Fields[0].Type, 0, false), g.toGolangType(kvt.Fields[1].Type, 0, false))
	}

	return casee.ToPascalCase(name)
}

func (g *generator) generateToGolangObject(srcvar, dstvar, rawtype string, indirections int) {
	if g.isListLike(rawtype) {
		g.generateListObjectToGolangSlice(srcvar, dstvar, rawtype)
	} else if g.isMapLike(rawtype) {
		g.generateMapObjectToGolangMap(srcvar, dstvar, rawtype)
	} else if _, ok := g.ctx.definedStruct[g.unwrapTypedef(rawtype)]; ok {
		// struct

		if indirections > 1 {
			panic("** field")
		}

		if indirections == 1 {
			g.printfln("%v = %v.toGoStruct()", dstvar, srcvar)
		} else {
			g.printfln("%v = *%v.toGoStruct()", dstvar, srcvar)
		}
	} else if _, ok := g.ctx.definedEnum[g.unwrapTypedef(rawtype)]; ok {
		g.printfln("%v = %v", dstvar, srcvar)
	} else if _, ok := g.ctx.definedInterface[g.unwrapTypedef(rawtype)]; ok {
		g.printfln("%v = %v", dstvar, srcvar)
	} else {
		// basic type
		t := g.toGolangType(rawtype, indirections, false)

		switch t {
		case "bool":
			fallthrough
		case "float64":
			fallthrough
		case "int32":
			fallthrough
		case "uint32":
			fallthrough
		case "uint64":
			fallthrough
		case "int64":
			fallthrough
		case "byte":
			fallthrough
		case "*byte": //??
			fallthrough
		case "FabricErrorCode":
			fallthrough
		case "uuid.UUID":
			fallthrough
		case "*uuid.UUID":
			g.printfln("%v = %v", dstvar, srcvar)
		case "interface{}":
			g.printfln("%v = fromUnsafePointer(%v)", dstvar, srcvar)
		case "string":
			g.printfln("%v = utf16PtrToString(%s)", dstvar, srcvar)
		case "time.Time":
			g.printfln("%v = %v.ToTime()", dstvar, srcvar)
		default:
			panic(fmt.Sprintf("unsupported generateToGolangObject type %v, raw %v", t, rawtype))
		}
	}
}
func (g *generator) generateListObjectToGolangSlice(srcvar, dstvar, listType string) {
	itemType := g.ctx.definedStruct[g.unwrapTypedef(listType)]
	countFieldName := itemType.Fields[0].Name
	itemFieldName := itemType.Fields[1].Name
	rawItemType := itemType.Fields[1].Type
	g.generatePointerToGolangSlice(dstvar, srcvar+"."+countFieldName, srcvar+"."+itemFieldName, rawItemType, 0)
}

func (g *generator) generatePointerToGolangSlice(dstSlice, countFieldName, itemFieldName, rawItemType string, indirections int) {
	golangItemTypeName := g.toGolangType(rawItemType, indirections, false)
	itemTypeName := g.toGolangType(rawItemType, indirections, true)
	// varuid := g.nextvarid()

	data := struct {
		GolangTypeName     string
		GolangItemTypeName string
		ItemTypeName       string
		CountFieldName     string
		ItemFieldName      string
		VarUid             int
		DstSlice           string
	}{
		GolangItemTypeName: golangItemTypeName,
		ItemTypeName:       itemTypeName,
		CountFieldName:     countFieldName,
		ItemFieldName:      itemFieldName,
		DstSlice:           dstSlice,
	}

	g.templateln(`{
		var innerlst []{{.ItemTypeName}}
		sliceCast(unsafe.Pointer(&innerlst), unsafe.Pointer({{.ItemFieldName}}), int({{.CountFieldName}}))

		for _, item := range innerlst {
			var tmpitem {{.GolangItemTypeName}}
	`, data)

	g.generateToGolangObject("item", "tmpitem", rawItemType, 0)

	g.templateln(`
			{{.DstSlice}} = append({{.DstSlice}}, tmpitem)
		}

	}`, data)
}

func (g *generator) generateMapObjectToGolangMap(srcvar, dstvar, mapType string) {
	// g.importpkg("reflect")

	golangTypeName := g.toGolangStructType(mapType, false)
	itemType := g.ctx.definedStruct[g.unwrapTypedef(mapType)]
	golangItemTypeName := g.toGolangType(itemType.Fields[1].Type, 0, false)

	kvType := g.ctx.definedStruct[g.unwrapTypedef(itemType.Fields[1].Type)]
	itemTypeName := g.toGolangType(itemType.Fields[1].Type, 0, true)
	varuid := g.nextvarid()

	data := struct {
		GolangTypeName     string
		GolangItemTypeName string
		ItemTypeName       string
		KeyTypeName        string // always string
		ValueTypeName      string
		Srcvar             string
		Dstvar             string
		CountFieldName     string
		ItemFieldName      string
		VarUid             int
	}{
		GolangTypeName:     golangTypeName,
		GolangItemTypeName: golangItemTypeName,
		ItemTypeName:       itemTypeName,
		KeyTypeName:        g.toGolangType(kvType.Fields[0].Type, 0, false),
		ValueTypeName:      g.toGolangType(kvType.Fields[1].Type, kvType.Fields[1].Indirections, false),
		Srcvar:             srcvar,
		Dstvar:             dstvar,
		CountFieldName:     itemType.Fields[0].Name,
		ItemFieldName:      itemType.Fields[1].Name,
		VarUid:             varuid,
	}

	g.templateln(`{
		var mapvar = make({{.GolangTypeName}})

		var innerlst []{{.ItemTypeName}}
		
		sliceCast(unsafe.Pointer(&innerlst), unsafe.Pointer({{.Srcvar}}.{{.ItemFieldName}}), int({{.Srcvar}}.{{.CountFieldName}}))

		for _, kv := range innerlst {
			var k {{.KeyTypeName}}
			var v {{.ValueTypeName}}
	`, data)

	// key
	g.generateToGolangObject(fmt.Sprintf("kv.%v", kvType.Fields[0].Name), "k", kvType.Fields[0].Type, 0)

	// val
	g.generateToGolangObject(fmt.Sprintf("kv.%v", kvType.Fields[1].Name), "v", kvType.Fields[1].Type, kvType.Fields[1].Indirections)

	if kvType.Fields[1].Indirections == 1 {
		g.printfln("mapvar[k] = *v")
	} else {
		g.printfln("mapvar[k] = v")
	}

	g.templateln(`
		}

		{{.Dstvar}} = mapvar
	}`, data)
}

func (g *generator) generateToInnerObject(srcvar, dstvar, rawtype string, indirections int) {
	if g.isListLike(rawtype) {
		g.generateSliceToInnerObject(srcvar, dstvar, rawtype)
	} else if g.isMapLike(rawtype) {
		g.generateMapToInnerObject(srcvar, dstvar, rawtype)
	} else if _, ok := g.ctx.definedStruct[g.unwrapTypedef(rawtype)]; ok {
		// struct

		if indirections > 1 {
			panic("** field")
		}

		if indirections == 1 {
			g.printfln("%v = %v.toInnerStruct()", dstvar, srcvar)
		} else {
			g.printfln("%v = *%v.toInnerStruct()", dstvar, srcvar)
		}
	} else if _, ok := g.ctx.definedEnum[g.unwrapTypedef(rawtype)]; ok {
		g.printfln("%v = %v", dstvar, srcvar)
	} else {
		// basic type
		t := g.toGolangType(rawtype, indirections, false)

		switch t {
		case "bool":
			fallthrough
		case "float64":
			fallthrough
		case "int32":
			fallthrough
		case "uint32":
			fallthrough
		case "uint64":
			fallthrough
		case "int64":
			fallthrough
		case "byte":
			fallthrough
		case "*byte": //??
			fallthrough
		case "FabricErrorCode":
			fallthrough
		case "uuid.UUID":
			g.printfln("%v = %v", dstvar, srcvar)
		case "interface{}":
			g.printfln("%v = toUnsafePointer(%v)", dstvar, srcvar)
		case "string":
			g.printfln("%v = utf16PtrFromString(%v)", dstvar, srcvar)
		case "time.Time":
			g.printfln("%v = timeToFiletime(%v)", dstvar, srcvar)
		default:
			panic(fmt.Sprintf("unsupported generateToInnerObject type %v, raw %v", t, rawtype))
		}
	}
}

func (g *generator) generateSliceToInnerObject(srcvar, dstvar, listType string) {
	innerTypeName := g.toGolangStructType(listType, true)

	uid := g.nextvarid()
	tmpvar := fmt.Sprintf("lst_%v", uid)
	g.printfln("%v := &%v{}", tmpvar, innerTypeName)

	itemType := g.ctx.definedStruct[g.unwrapTypedef(listType)]
	countFieldName := tmpvar + "." + itemType.Fields[0].Name
	itemFieldName := tmpvar + "." + itemType.Fields[1].Name
	rawItemType := itemType.Fields[1].Type

	g.generateSliceToPointer(srcvar, countFieldName, itemFieldName, rawItemType, 0)
	g.printfln("%v = %v", dstvar, tmpvar)
}

func (g *generator) generateSliceToPointer(srcSlice, countFieldName, itemFieldName, rawItemType string, indirections int) {
	itemTypeName := g.toGolangType(rawItemType, indirections, true)

	data := struct {
		ItemTypeName   string
		Srcvar         string
		CountFieldName string
		ItemFieldName  string
	}{
		ItemTypeName:   itemTypeName,
		Srcvar:         srcSlice,
		CountFieldName: countFieldName,
		ItemFieldName:  itemFieldName,
	}

	g.templateln(`{

		var tmp []{{.ItemTypeName}}

		for _, item := range {{.Srcvar}} {
			var tmpitem {{.ItemTypeName}}
	`, data)

	g.generateToInnerObject("item", "tmpitem", rawItemType, 0)

	g.templateln(`
			tmp = append(tmp, tmpitem)
		}

		{{.CountFieldName}} = uint32(len(tmp))
		if len(tmp) > 0 {
			{{.ItemFieldName}} = &tmp[0]
		}

	}`, data)
}

func (g *generator) generateMapToInnerObject(srcvar, dstvar, mapType string) {
	innerTypeName := g.toGolangStructType(mapType, true)
	itemType := g.ctx.definedStruct[g.unwrapTypedef(mapType)]
	kvType := g.ctx.definedStruct[g.unwrapTypedef(itemType.Fields[1].Type)]
	itemTypeName := g.toGolangType(itemType.Fields[1].Type, 0, true)
	varuid := g.nextvarid()

	data := struct {
		InnerTypeName  string
		ItemTypeName   string
		KVTypeName     string
		Srcvar         string
		Dstvar         string
		CountFieldName string
		ItemFieldName  string
		VarUid         int
	}{
		InnerTypeName:  innerTypeName,
		ItemTypeName:   itemTypeName,
		KVTypeName:     g.toGolangType(kvType.Name, 0, true),
		Srcvar:         srcvar,
		Dstvar:         dstvar,
		CountFieldName: itemType.Fields[0].Name,
		ItemFieldName:  itemType.Fields[1].Name,
		VarUid:         varuid,
	}

	// TODO more than string
	g.templateln(`{
		mapobj := &{{.InnerTypeName}}{}

		var tmp []{{.ItemTypeName}}

		for k, v := range {{.Srcvar}} {
			kv := {{.KVTypeName}}{}
	`, data)

	// key
	g.generateToInnerObject("k", fmt.Sprintf("kv.%v", kvType.Fields[0].Name), kvType.Fields[0].Type, 0)

	// val
	g.generateToInnerObject("v", fmt.Sprintf("kv.%v", kvType.Fields[1].Name), kvType.Fields[1].Type, kvType.Fields[1].Indirections)

	g.templateln(`
		tmp = append(tmp, kv)
		}

		mapobj.{{.CountFieldName}} = uint32(len(tmp))
		if len(tmp) > 0 {
			mapobj.{{.ItemFieldName}} = &tmp[0]
		}

		{{.Dstvar}} = mapobj
	}`, data)
}
