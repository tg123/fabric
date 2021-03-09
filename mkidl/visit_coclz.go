package main

import (
	"fmt"
	"strings"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/jd3nn1s/gomidl/scanner"
	"github.com/pinzolo/casee"
)

func (g *generator) generateCoClz(n *ast.CoClassNode) {
	hubName := casee.ToCamelCase(n.Name)
	clzName := n.Name

	g.printfln("type %vComHub struct{", hubName)
	for _, ifc := range n.Interfaces {
		g.printfln("%v *%v", strings.TrimPrefix(ifc.Name, "I"), g.goInterfaceName(ifc.Name))
	}
	g.printfln("}")

	// init
	g.printfln("func (v *%vComHub) init(createComObject comCreator) {", hubName)

next:
	for _, ifc := range n.Interfaces {
		ifc = g.ctx.definedInterface[ifc.Name]
		for _, attr := range ifc.Attributes {
			if attr.Type == scanner.UUID {
				g.printfln(`createComObject("{%v}", unsafe.Pointer(&v.%v))`, attr.Val, strings.TrimPrefix(ifc.Name, "I"))
				continue next
			}
		}
	}
	g.printfln("}")

	// methods
	for _, ifc := range n.Interfaces {
		ifc = g.ctx.definedInterface[ifc.Name]
		hubFieldName := strings.TrimPrefix(ifc.Name, "I")

		for i, m := range ifc.Methods {
			if methodblacklist[fmt.Sprintf("%v.%v", ifc.Name, m.Name)] {
				continue
			}

			if strings.HasPrefix(m.Name, "Begin") {
				continue
			}

			var params []string
			var method string

			// TODO support EndGetNodeList2
			if strings.HasPrefix(m.Name, "End") {
				begin := ifc.Methods[i-1]
				end := m
				_, params, method = g.generateAsyncCallSig(clzName, begin, end, true)
			} else {
				params, method = g.generateMethodSig(clzName, m, true)
			}

			g.printfln(`if v.hub.%v == nil {
				err = errComNotImpl
				return
			}`, hubFieldName)
			g.printfln("return v.hub.%v.%v(%v)", hubFieldName, method, strings.Join(params, ","))
			g.printfln("}")
		}
	}

}
