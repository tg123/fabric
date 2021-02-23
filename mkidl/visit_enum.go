package main

import (
	"fmt"
	"strings"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/pinzolo/casee"
)

func goEnumName(id string) string {
	// hack for FABRIC_ERROR_CODE
	if id == "FABRIC_ERROR_CODE" {
		return "FabricErrorCode"
	} else if strings.HasPrefix(id, "FABRIC_E") {
		id = "FABRIC_ERROR" + strings.TrimPrefix(id, "FABRIC_E")
	}

	return casee.ToPascalCase(id)
}

// val of enum node is raw string expr, replace any ident with new go name
func (g *generator) parseEnumExpr(expr, casttype string) string {
	if expr == "" {
		return ""
	}

	// numbers
	if strings.HasPrefix(expr, "0") {
		return expr
	}

	if strings.HasPrefix(expr, "-") {
		return expr
	}

	// fast return for single ident
	if refEnumType, ok := g.ctx.definedEnumValueType[expr]; ok {
		if refEnumType != casttype {
			return fmt.Sprintf("%v(%v)", goEnumName(casttype), goEnumName(expr))
		}
	}

	// convert all ident in expr
	e := expr
	for refEnumValue, refEnumType := range g.ctx.definedEnumValueType {
		if strings.Contains(e, refEnumValue) {

			v := goEnumName(refEnumValue)

			if refEnumType != casttype {
				v = fmt.Sprintf("%v(%v)", goEnumName(casttype), v)
			}

			e = strings.ReplaceAll(e, refEnumValue, v)
		}
	}

	return e
}

func (g *generator) visitEnum(n *ast.EnumNode) {

	t := "int32"
	if n.Name == "FABRIC_ERROR_CODE" {
		t = "int64"
	}

	enumType := goEnumName(n.Name)

	g.printfln("type %s %v", enumType, t)
	g.printfln("const (")

	for i, v := range n.Values {
		enumName := goEnumName(v.Name)
		enumVal := g.parseEnumExpr(v.Val, n.Name)

		// convert to iota
		if (enumVal == "") && i == 0 {
			enumVal = "iota"
		} else if i < len(n.Values)-1 {
			if n.Values[i+1].Val == "" {
				if enumVal != "" {
					enumVal = fmt.Sprintf("iota - %v + %v", i, enumVal)
				}

			}
		}

		if enumVal == "" {
			g.printfln("%v", enumName)
		} else {
			g.printfln("%s %s = %s", enumName, enumType, enumVal)
		}
	}
	g.printfln(")")
}
