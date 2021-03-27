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

var interfaceBlackList = map[string]bool{
	// TODO goproxy to support inherit IFabricClientConnectionEventHandler
	"IFabricClientConnectionEventHandler2": true,
}

type goproxyProperty struct {
	Name   string
	Type   string
	NoCtor bool
}

var goproxyProperties = map[string][]goproxyProperty{
	"IFabricServiceNotificationEventHandler": {
		{
			Name: "client",
			Type: "*FabricClient",
		},
	},
	"IFabricClientConnectionEventHandler": {
		{
			Name: "client",
			Type: "*FabricClient",
		},
	},
	"IFabricAsyncOperationCallback": {
		{
			Name: "callback",
			Type: "func(ctx *comFabricAsyncOperationContext)",
		},
	},
	"IFabricAsyncOperationContext": {
		{
			Name: "nativeCallback",
			Type: "*comFabricAsyncOperationCallback",
		},
		{
			Name:   "result",
			Type:   "interface{}",
			NoCtor: true,
		},
		{
			Name:   "resultHResult",
			Type:   "uintptr",
			NoCtor: true,
		},
		{
			Name:   "lock",
			Type:   "sync.Mutex",
			NoCtor: true,
		},
		{
			Name: "goctx",
			Type: "context.Context",
		},
		{
			Name: "cancel",
			Type: "context.CancelFunc",
		},
	},
	"IFabricStringResult": {
		{
			Name: "result",
			Type: "string",
		},
	},
	"IFabricStatelessServiceFactory": {
		{
			Name: "builder",
			Type: "func(ServiceContext) (StatelessServiceInstance, error)",
		},
	},
	"IFabricStatelessServiceInstance": {
		{
			Name: "instance",
			Type: "StatelessServiceInstance",
		},
	},
}

func (g *generator) goInterfaceName(name string) string {

	// if g.ctx.publicInterfaces[name] {
	// 	return fmt.Sprintf("Com%v", strings.TrimPrefix(name, "I"))
	// }

	return fmt.Sprintf("com%v", strings.TrimPrefix(name, "I"))
}

// TODO support those method
var methodBlackList = map[string]bool{

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
	"IFabricPropertyManagementClient.EndSubmitPropertyBatch":    true,
	"IFabricNameEnumerationResult.GetNames":                     true,

	// []
	"IFabricApplicationUpgradeProgressResult.GetUpgradeDomains":        true,
	"IFabricApplicationUpgradeProgressResult.GetChangedUpgradeDomains": true,
	"IFabricStringListResult.GetStrings":                               true,

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

func (g *generator) generateMethodSig(receiver string, m *ast.MethodNode, forceHidden bool) ([]string, string) {
	methodName := goMethodName(m.Name)

	if forceHidden {
		g.printfln("func (v *%s) %s(", receiver, casee.ToCamelCase(methodName))
	} else {
		g.printfln("func (v *%s) %s(", receiver, methodName)
	}
	var rt []string
	var paramNames []string

	if m.ReturnType.Type != "HRESULT" {
		rt = append(rt, fmt.Sprintf("rt %v", g.toGolangType(m.ReturnType.Type, m.ReturnType.Indirections, false)))
	}

	for i, p := range m.Params {
		paramName := p.Name
		if paramName == "" {
			paramName = fmt.Sprintf("param_%v", i)
		}

		if isOutParam(p) {
			rt = append(rt, fmt.Sprintf("%v %v", paramName, g.toGolangType(p.Type, p.Indirections-1, false)))
		} else {
			g.printfln("%v %v,", paramName, g.toGolangType(p.Type, p.Indirections, false))
			paramNames = append(paramNames, paramName)
		}
	}
	rt = append(rt, "err error")
	g.printfln(") (%v) {", strings.Join(rt, ","))

	return paramNames, methodName
}

func (g *generator) generateMethods(n *ast.InterfaceNode) {
	interfaceName := g.goInterfaceName(n.Name)

	for _, m := range n.Methods {

		// TODO fix blacklist
		if methodBlackList[fmt.Sprintf("%v.%v", n.Name, m.Name)] {
			continue
		}

		g.generateMethodSig(interfaceName, m, false)
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
				case "*ole.GUID":
					syscallParams = append(syscallParams, fmt.Sprintf("uintptr(unsafe.Pointer(%s))", paramName))
				case "ole.GUID":
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

}

func (g *generator) generateComStub(n *ast.InterfaceNode) {
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
			{{if .HasGoProxy}} proxy {{.Name | ToCamelCase }}GoProxy {{end}}
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
		HasGoProxy  bool
	}{
		Name:        interfaceName,
		InnerName:   casee.ToCamelCase(interfaceName),
		Parent:      pn,
		InnerParent: casee.ToCamelCase(pn),
		Methods:     n.Methods,
		HasGoProxy:  goproxyProperties[n.Name] != nil,
	})

	g.generateMethods(n)
}

func (g *generator) generateGoProxy(n *ast.InterfaceNode) {
	g.importpkg("syscall")
	interfaceName := g.goInterfaceName(n.Name)
	props := goproxyProperties[n.Name]

	for _, p := range props {
		p := strings.Split(p.Type, ".")
		if len(p) == 2 {
			g.importpkg(p[0])
		}
	}

	g.templateln(`
		type {{.Name | ToCamelCase }}GoProxy struct {
			unknownref *goIUnknown
			{{ range .Properties }} {{.Name}} {{.Type}}
			{{ end }} 
		}

		func new{{ .Name | ToPascalCase }}( 
			{{ range .Properties }} {{if not .NoCtor }} {{.Name}} {{.Type}}, {{end}} 
			{{end}}) *{{.Name}} {
			com := &{{.Name}}{}
			*(**{{.InnerName}}Vtbl)(unsafe.Pointer(com)) = &{{.InnerName}}Vtbl{}
			vtbl := com.vtable()
			com.proxy.unknownref = attachIUnknown("{{"{"}}{{ .IID }}{{"}"}}", &vtbl.IUnknownVtbl)
			{{ range .Methods }} vtbl.{{.Name}} = syscall.NewCallback(com.proxy.{{.Name | ToPascalCase }})
			{{ end }} 
			{{ range .Properties }}  {{if not .NoCtor }} com.proxy.{{.Name}} = {{.Name}} {{ end }} 
			{{ end }} 
			com.proxy.init()
			return com
		}
	`, struct {
		Name       string
		InnerName  string
		Methods    []*ast.MethodNode
		Properties []goproxyProperty
		IID        string
	}{
		Name:       interfaceName,
		InnerName:  casee.ToCamelCase(interfaceName),
		Methods:    n.Methods,
		Properties: props,
		IID:        strings.ToUpper(findIID(n)),
	})

	g.printfln("/*")
	for _, m := range n.Methods {
		g.printfln("func (v *%sGoProxy) %s(", interfaceName, casee.ToPascalCase(m.Name))
		g.printfln("_ *ole.IUnknown,")
		for i, p := range m.Params {
			paramName := p.Name
			if paramName == "" {
				paramName = fmt.Sprintf("param_%v", i)
			}
			g.printfln("%v %v,", paramName, g.toGolangType(p.Type, p.Indirections, true))
		}
		g.printfln(") uintptr { return 0}")
	}
	g.printfln("*/")
}

func (g *generator) visitInterface(n *ast.InterfaceNode) {
	if _, ok := interfaceBlackList[n.Name]; ok {
		return
	}

	if _, ok := goproxyProperties[n.Name]; ok {
		g.generateGoProxy(n)
	}
	g.generateComStub(n)
}
