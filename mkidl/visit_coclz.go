package main

import (
	"fmt"
	"strings"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/jd3nn1s/gomidl/scanner"
	"github.com/pinzolo/casee"
)

func (g *generator) generateAsyncCallSig(receiver string, begin, end *ast.MethodNode) (*ast.InterfaceNode, []string, string) {
	var paramNames []string
	methodName := strings.TrimPrefix(begin.Name, "Begin")
	g.printfln("func (v *%s) %s(", receiver, methodName)
	g.importpkg("context")
	g.printfln("ctx context.Context,")
	paramNames = append(paramNames, "ctx")

	for i, p := range begin.Params {

		if p.Type == "IFabricAsyncOperationCallback" {
			continue
		}

		if p.Type == "IFabricAsyncOperationContext" {
			continue
		}

		// replace by context.Context
		if p.Name == "timeoutMilliseconds" {
			continue
		}

		paramName := p.Name
		if paramName == "" {
			paramName = fmt.Sprintf("param_%v", i)
		}
		g.printfln("%v %v,", paramName, g.toGolangType(p.Type, p.Indirections, false))
		paramNames = append(paramNames, paramName)
	}

	var rt []string

	var asyncrt *ast.InterfaceNode
	for _, p := range end.Params {
		for _, a := range p.Attributes {
			if a.Type == scanner.RETVAL {

				asyncrt = g.ctx.definedInterface[g.unwrapTypedef(p.Type)]

				if asyncrt == nil {
					rt = append(rt, fmt.Sprintf("result_%v %v", 0, g.toGolangType(p.Type, p.Indirections-1, false)))
					break
				}

				for i, m := range g.ctx.definedInterface[g.unwrapTypedef(p.Type)].Methods {
					rt = append(rt, fmt.Sprintf("result_%v %v", i, g.toGolangType(m.ReturnType.Type, m.ReturnType.Indirections, false)))
				}

			}
		}
	}

	rt = append(rt, "err error")
	g.printfln(") (%v) {", strings.Join(rt, ","))

	return asyncrt, paramNames, methodName
}

// TODO moved from interface, should clean up for coclz
func (g *generator) generateAsyncCall(hubFieldName string, asyncrt *ast.InterfaceNode, n *ast.InterfaceNode, begin, end *ast.MethodNode) {
	// asyncrt, _, _ := g.generateAsyncCallSig(g.goInterfaceName(n.Name), begin, end, true)
	g.printfln(`if v.hub.%v == nil {
		err = errComNotImpl
		return
	}`, hubFieldName)

	g.printfln(`ch := make(chan error, 1)
		defer close(ch)
		callback := newComFabricAsyncOperationCallback(func(sfctx *comFabricAsyncOperationContext) {
	`)

	var rt []string
next:
	for i, p := range end.Params {
		if p.Type == "IFabricAsyncOperationContext" {
			continue
		}

		if isOutParam(p) {
			rt = append(rt, fmt.Sprintf("rt_%v", i))
			continue next
		}

		panic("non out end param")
	}

	rt = append(rt, "err")
	g.printfln("%v := v.hub.%v.%v(sfctx)", strings.Join(rt, ","), hubFieldName, casee.ToCamelCase(end.Name))
	g.printfln(`
	if err != nil {
		ch <- err
		return
	}`)

	if len(rt) > 3 {
		panic(fmt.Sprintf("too many end rt %v", end.Name))
	}

	if asyncrt != nil {
		g.printfln(`defer releaseComObject(&%v.IUnknown)`, rt[0])
		for i, m := range asyncrt.Methods {
			g.printfln("result_%v, err = %v.%v()", i, rt[0], casee.ToPascalCase(m.Name))

			g.printfln(`if err != nil {
				ch <- err
				return
			}`)
		}

	} else {
		g.printfln(`if err != nil {
			ch <- err
			return
		}`)

		if len(rt) == 2 {
			g.printfln("result_%v = %v", 0, rt[0])
		}
	}

	g.printfln(`ch <- nil
			})`)

	g.printfln(`timeout := toTimeout(ctx, v)`)

	var beginrt []string
	for _, p := range begin.Params {
		if p.Type == "IFabricAsyncOperationContext" {
			beginrt = append(beginrt, "sfctx")
		} else if isOutParam(p) {
			beginrt = append(beginrt, "_")
		}
	}

	beginrt = append(beginrt, "err")

	g.printfln("%v := v.hub.%v.%v(", strings.Join(beginrt, ","), hubFieldName, casee.ToCamelCase(begin.Name))
	for _, p := range begin.Params {
		if p.Type == "IFabricAsyncOperationContext" {
			continue
		}

		if isOutParam(p) {
			continue
		}

		if p.Type == "IFabricAsyncOperationCallback" {
			g.printfln("callback,")
			continue
		}

		if p.Name == "timeoutMilliseconds" {
			g.printfln("uint32(timeout.Milliseconds()),")
			continue
		}

		g.printfln("%s, ", p.Name)
	}
	g.printfln(") ")

	g.importpkg("runtime")
	g.printfln(`
	if err != nil {
		return
	}

	err = waitch(ctx, ch, sfctx, timeout)
	runtime.KeepAlive(callback)
	return
	`)

	g.printfln("}")
}

func findIID(ifc *ast.InterfaceNode) string {
	for _, attr := range ifc.Attributes {
		if attr.Type == scanner.UUID {
			return attr.Val
		}
	}

	return ""
}

var hubHiddenList = map[string]bool{
	"IFabricClientSettings2.GetSettings":                                              true,
	"IFabricRuntime.RegisterStatelessServiceFactory":                                  true,
	"IFabricRuntime.RegisterStatefulServiceFactory":                                   true,
	"IFabricRuntime.CreateServiceGroupFactoryBuilder":                                 true,
	"IFabricRuntime.RegisterServiceGroupFactory":                                      true,
	"IFabricCodePackageActivationContext.RegisterCodePackageChangeHandler":            true,
	"IFabricCodePackageActivationContext.UnregisterCodePackageChangeHandler":          true,
	"IFabricCodePackageActivationContext.RegisterConfigurationPackageChangeHandler":   true,
	"IFabricCodePackageActivationContext.UnregisterConfigurationPackageChangeHandler": true,
	"IFabricCodePackageActivationContext.RegisterDataPackageChangeHandler":            true,
	"IFabricCodePackageActivationContext.UnregisterDataPackageChangeHandler":          true,
}

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

	for _, ifc := range n.Interfaces {
		ifc = g.ctx.definedInterface[ifc.Name]
		hubFieldName := strings.TrimPrefix(ifc.Name, "I")

		iid := findIID(ifc)

		if iid != "" {
			g.importpkg("unsafe")
			g.printfln(`createComObject("{%v}", unsafe.Pointer(&v.%v))`, strings.ToUpper(iid), hubFieldName)
		}
	}
	g.printfln("}")

	g.printfln("func (v *%vComHub) Close() error {", hubName)
	for _, ifc := range n.Interfaces {
		ifc = g.ctx.definedInterface[ifc.Name]
		hubFieldName := strings.TrimPrefix(ifc.Name, "I")

		if findIID(ifc) != "" {
			g.printfln(`if v.%v != nil {`, hubFieldName)
			g.printfln(`releaseComObject(&v.%v.IUnknown)`, hubFieldName)
			g.printfln("}")
		}
		g.printfln("return nil")
	}
	g.printfln("}")

	// methods
	for _, ifc := range n.Interfaces {
		ifc = g.ctx.definedInterface[ifc.Name]
		hubFieldName := strings.TrimPrefix(ifc.Name, "I")

		for i, m := range ifc.Methods {
			if methodBlackList[fmt.Sprintf("%v.%v", ifc.Name, m.Name)] {
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
				var asyncrt *ast.InterfaceNode

				// TODO unify hidden
				asyncrt, params, method = g.generateAsyncCallSig(clzName, begin, end)

				g.generateAsyncCall(hubFieldName, asyncrt, ifc, begin, end)

			} else {
				params, method = g.generateMethodSig(clzName, m, hubHiddenList[fmt.Sprintf("%v.%v", ifc.Name, m.Name)])

				g.printfln(`if v.hub.%v == nil {
					err = errComNotImpl
					return
				}`, hubFieldName)
				g.printfln("return v.hub.%v.%v(%v)", hubFieldName, method, strings.Join(params, ","))
				g.printfln("}")
			}

		}
	}

}
