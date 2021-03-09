package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/jd3nn1s/gomidl/scanner"
	"github.com/pinzolo/casee"
)

func (g *generator) goInterfaceName(name string) string {

	if g.ctx.publicReturnedInterfaces[name] {
		return fmt.Sprintf("Com%v", strings.TrimPrefix(name, "I"))
	}

	return fmt.Sprintf("com%v", strings.TrimPrefix(name, "I"))
}

// TODO support those method
var methodblacklist = map[string]bool{

	// weird
	"IFabricQueryClient6.EndGetNodeList2":        true,
	"IFabricQueryClient6.EndGetApplicationList2": true,
	"IFabricQueryClient6.EndGetServiceList2":     true,
	"IFabricQueryClient6.EndGetPartitionList2":   true,
	"IFabricQueryClient6.EndGetReplicaList2":     true,

	// with sync
	"IFabricRuntime.EndRegisterStatelessServiceFactory": true,
	"IFabricRuntime.EndRegisterStatefulServiceFactory":  true,
	"IFabricRuntime.EndRegisterServiceGroupFactory":     true,

	// fix IFabricPropertyEnumerationResult
	"IFabricPropertyManagementClient.EndEnumerateProperties":    true,
	"IFabricPropertyManagementClient.EndGetProperty":            true,
	"IFabricPropertyManagementClient.EndEnumerateSubNames":      true,
	"IFabricServiceManagementClient.EndResolveServicePartition": true,
	"IFabricNameEnumerationResult.GetNames":                     true,

	// []
	"IFabricApplicationUpgradeProgressResult.GetUpgradeDomains":        true,
	"IFabricApplicationUpgradeProgressResult.GetChangedUpgradeDomains": true,
	"IFabricStringListResult.GetStrings":                               true,

	// generator bug skip
	"IFabricPropertyManagementClient.EndSubmitPropertyBatch": true,

	"IFabricStatelessServiceInstance.EndOpen": true,
	"IFabricStatefulServiceReplica.EndOpen":   true,
	"IFabricReplicator.EndOpen":               true,
	"IFabricOperationStream.EndGetOperation":  true,
	"IFabricOperationDataStream.EndGetNext":   true,

	// TODO Name Duplicate
	"IFabricFaultManagementClient.EndRestartDeployedCodePackage": true,
	"IFabricFaultManagementClient.EndStartNode":                  true,
	"IFabricFaultManagementClient.EndStopNode":                   true,
	"IFabricFaultManagementClient.EndRestartNode":                true,
}

func goMethodName(m string) string {
	if strings.HasPrefix(m, "Begin") || strings.HasPrefix(m, "End") {
		return casee.ToCamelCase(m)
	}
	return casee.ToPascalCase(m)

}

func isOutParam(p *ast.ParamNode) bool {
	for _, a := range p.Attributes {
		if a.Type == scanner.OUT {
			return true
		}
	}

	return false
}

func (g *generator) generateAsyncCallSig(receiver string, begin, end *ast.MethodNode, namedRt bool) (*ast.InterfaceNode, []string, string) {
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
					if namedRt {
						rt = append(rt, fmt.Sprintf("result_%v %v", 0, g.toGolangType(p.Type, p.Indirections-1, false)))
					} else {
						rt = append(rt, fmt.Sprintf("%v", g.toGolangType(p.Type, p.Indirections-1, false)))
					}
					break
				}

				for i, m := range g.ctx.definedInterface[g.unwrapTypedef(p.Type)].Methods {
					if namedRt {
						rt = append(rt, fmt.Sprintf("result_%v %v", i, g.toGolangType(m.ReturnType.Type, m.ReturnType.Indirections, false)))

					} else {
						rt = append(rt, fmt.Sprintf("%v", g.toGolangType(m.ReturnType.Type, m.ReturnType.Indirections, false)))

					}
				}

			}
		}
	}

	if namedRt {
		rt = append(rt, "err error")
	} else {
		rt = append(rt, "error")
	}
	g.printfln(") (%v) {", strings.Join(rt, ","))

	return asyncrt, paramNames, methodName
}

func (g *generator) generateAsyncCall(n *ast.InterfaceNode, begin, end *ast.MethodNode) {
	asyncrt, _, _ := g.generateAsyncCallSig(g.goInterfaceName(n.Name), begin, end, true)

	g.printfln(`ch := make(chan error, 1)
	defer close(ch)
	callback := newFabricAsyncOperationCallback(func(sfctx *comIFabricAsyncOperationContext) {
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
	g.printfln("%v := v.%v(sfctx)", strings.Join(rt, ","), casee.ToCamelCase(end.Name))
	g.printfln(`
	if err != nil {
		ch <- err
		return
	}`)

	if len(rt) > 3 {
		panic(fmt.Sprintf("too many end rt %v", end.Name))
	}

	if asyncrt != nil {
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

	g.printfln(`timeout := toTimeout(ctx)`)

	var beginrt []string
	for _, p := range begin.Params {
		if p.Type == "IFabricAsyncOperationContext" {
			beginrt = append(beginrt, "sfctx")
		} else if isOutParam(p) {
			beginrt = append(beginrt, "_")
		}
	}

	beginrt = append(beginrt, "err")

	g.printfln("%v := v.%v(", strings.Join(beginrt, ","), casee.ToCamelCase(begin.Name))
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

	g.printfln(`
	if err != nil {
		return
	}

	err = waitch(ctx, ch, sfctx, timeout)
	return
	`)

	g.printfln("}")
}

func (g *generator) generateMethodSig(receiver string, m *ast.MethodNode, namedRt bool) ([]string, string) {
	methodName := goMethodName(m.Name)
	g.printfln("func (v *%s) %s(", receiver, methodName)

	var rt []string
	var paramNames []string

	if m.ReturnType.Type != "HRESULT" {
		if namedRt {
			rt = append(rt, fmt.Sprintf("rt %v", g.toGolangType(m.ReturnType.Type, m.ReturnType.Indirections, false)))
		} else {
			rt = append(rt, fmt.Sprintf("%v", g.toGolangType(m.ReturnType.Type, m.ReturnType.Indirections, false)))
		}
	}

	for i, p := range m.Params {
		paramName := p.Name
		if paramName == "" {
			paramName = fmt.Sprintf("param_%v", i)
		}

		if isOutParam(p) {
			if namedRt {
				rt = append(rt, fmt.Sprintf("%v %v", paramName, g.toGolangType(p.Type, p.Indirections-1, false)))
			} else {
				rt = append(rt, fmt.Sprintf("%v", g.toGolangType(p.Type, p.Indirections-1, false)))
			}
		} else {
			g.printfln("%v %v,", paramName, g.toGolangType(p.Type, p.Indirections, false))
			paramNames = append(paramNames, paramName)
		}
	}
	if namedRt {
		rt = append(rt, "err error")
	} else {
		rt = append(rt, "error")
	}
	g.printfln(") (%v) {", strings.Join(rt, ","))

	return paramNames, methodName
}

func (g *generator) generateMethods(n *ast.InterfaceNode) {
	interfaceName := g.goInterfaceName(n.Name)

	for _, m := range n.Methods {

		// TODO fix blacklist
		if methodblacklist[fmt.Sprintf("%v.%v", n.Name, m.Name)] {
			continue
		}

		g.generateMethodSig(interfaceName, m, true)
		syscallParams := make([]string, 0, len(m.Params))

		// TODO dup code, but i did not find better way to make it more clear
		for i, p := range m.Params {
			paramName := p.Name
			if paramName == "" {
				paramName = fmt.Sprintf("param_%v", i)
			}

			if isOutParam(p) {
				// fmt.Printf("%v.%v\n", n.Name, m.Name)
				t := g.toGolangType(p.Type, p.Indirections-1, true)
				g.printfln("var p_%v %v", i, t)
				g.printfln(`defer func(){`)

				g.generateToGolangObject(fmt.Sprintf("p_%v", i), paramName, p.Type, p.Indirections-1)
				g.printfln(`}()`)
				syscallParams = append(syscallParams, fmt.Sprintf("uintptr(unsafe.Pointer(&p_%v))", i))
			} else {
				t := g.toGolangType(p.Type, p.Indirections, true)
				switch t {
				case "bool":
					g.printfln("p_%v := 0", i)
					g.printfln("if %v {", paramName)
					g.printfln("p_%v = 1", i)
					g.printfln("}")
					syscallParams = append(syscallParams, fmt.Sprintf("uintptr(p_%v)", i))
				case "windows.GUID":
					syscallParams = append(syscallParams, fmt.Sprintf("uintptr(unsafe.Pointer(&%s))", paramName))
				case "unsafe.Pointer": // interface{}
					syscallParams = append(syscallParams, fmt.Sprintf("uintptr(toUnsafePointer(%v))", paramName))
				default:

					if _, ok := g.ctx.definedStruct[g.unwrapTypedef(p.Type)]; ok || t == "*uint16" {
						// string or obj
						g.printfln("var p_%v %v", i, t)
						g.generateToInnerObject(paramName, fmt.Sprintf("p_%v", i), p.Type, p.Indirections)
						syscallParams = append(syscallParams, fmt.Sprintf("uintptr(unsafe.Pointer(p_%v))", i))
					} else if p.Indirections > 0 {
						// pointer
						syscallParams = append(syscallParams, fmt.Sprintf("uintptr(unsafe.Pointer(%v))", paramName))
					} else {
						// all others
						syscallParams = append(syscallParams, fmt.Sprintf("uintptr(%v)", paramName))
					}

				}

			}
		}

		numSyscallParams := len(m.Params) + 1
		numSyscallRequired := int(math.Ceil(float64(numSyscallParams)/3)) * 3

		syscallFunc := "Syscall"
		if numSyscallRequired > 18 {
			panic("more params than supported by syscall")
		} else if numSyscallRequired > 3 {
			syscallFunc += strconv.Itoa(numSyscallRequired)
		}

		actualSyscallParamLen := len(syscallParams)
		for i := len(syscallParams); i < numSyscallRequired-1; i++ {
			syscallParams = append(syscallParams, "0")
		}

		g.importpkg("unsafe")
		g.importpkg("syscall")

		g.printfln("hr, _, err1 := syscall.%s(", syscallFunc)
		g.printfln("v.vtable().%s,", m.Name)
		g.printfln("%d,", actualSyscallParamLen+1)
		g.printfln("uintptr(unsafe.Pointer(v)),")

		for _, param := range syscallParams {
			g.printfln("%v,", param)
		}

		g.printfln(")")

		rtype := g.toGolangType(m.ReturnType.Type, m.ReturnType.Indirections, true)
		if m.ReturnType.Type == "HRESULT" {
			// covert to error
			g.printfln(`if hr != 0 { 
				err = errno(hr, err1)
				return 
			}`)
		} else if _, ok := g.ctx.definedEnum[m.ReturnType.Type]; ok || rtype == "uint32" {
			// cast to enum or uint32 from unitptr
			g.printfln(`_ = err1
						rt = %v(hr)`, rtype)
		} else if rtype == "bool" {
			g.printfln(` _ = err1
			rt = hr != 0`)
		} else {

			// convert to obj
			g.printfln(`if hr == 0 { 
				err = err1
				return 
			}
			
			tmp := (%v)(unsafe.Pointer(hr))
			`, rtype)

			if g.ctx.definedInterface[m.ReturnType.Type] != nil {
				g.printfln("rt = tmp")
			} else {
				g.generateToGolangObject("tmp", "rt", m.ReturnType.Type, m.ReturnType.Indirections)
			}
		}

		g.printfln("return }")

	}

	for i := 1; i < len(n.Methods); i++ {
		begin := n.Methods[i-1]
		end := n.Methods[i]

		// TODO support EndGetNodeList2

		if strings.HasPrefix(begin.Name, "Begin") && strings.HasPrefix(end.Name, "End") {

			if methodblacklist[fmt.Sprintf("%v.%v", n.Name, begin.Name)] {
				continue
			}

			if methodblacklist[fmt.Sprintf("%v.%v", n.Name, end.Name)] {
				continue
			}

			g.generateAsyncCall(n, begin, end)
		}
	}
}

func (g *generator) visitInterface(n *ast.InterfaceNode) {
	if _, ok := basicTypeMap[n.Name]; ok {
		return
	}

	interfaceName := g.goInterfaceName(n.Name)

	// generate vtable
	pn := n.ParentName
	if _, ok := g.ctx.definedInterface[pn]; ok {
		pn = g.goInterfaceName(pn)
	}
	if pn == "IUnknown" || pn == "" {
		g.importpkg("ole")
		pn = "ole.IUnknown"
	}

	g.importpkg("unsafe")
	g.templateln(`
		type {{.Name}} struct {
			{{.Parent}}
		}

		type {{.InnerName}}Vtbl struct {
			{{.InnerParent}}Vtbl
			{{ range .Methods }} {{.Name}} uintptr
			{{ end }} 
		}

		func (v *{{.Name}}) vtable() *{{.InnerName}}Vtbl {
			return (*{{.InnerName}}Vtbl)(unsafe.Pointer(v.RawVTable))
		}

	`, struct {
		Name        string
		InnerName   string
		Parent      string
		InnerParent string
		Methods     []*ast.MethodNode
	}{
		Name:        interfaceName,
		InnerName:   casee.ToCamelCase(interfaceName),
		Parent:      pn,
		InnerParent: casee.ToCamelCase(pn),
		Methods:     n.Methods,
	})

	g.generateMethods(n)
}
