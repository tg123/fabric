package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/jd3nn1s/gomidl/parser"
	"github.com/pinzolo/casee"
)

type defContext struct {
	definedIdls map[string]interface{}

	// for forward lookup
	definedEnum          map[string]*ast.EnumNode
	definedEnumValueType map[string]string
	definedStruct        map[string]*ast.StructNode
	definedInterface     map[string]*ast.InterfaceNode
	definedTypedef       map[string]*ast.TypedefNode

	publicReturnedInterfaces map[string]bool
}

func defineIdent(c *defContext, nodes []interface{}) {
	for _, n := range nodes {
		switch v := n.(type) {
		// case *ast.ImportNode:
		// case *ast.ConstdefNode:
		case *ast.StructNode:
			c.definedStruct[v.Name] = v
		case *ast.EnumNode:
			c.definedEnum[v.Name] = v

			for _, e := range v.Values {
				c.definedEnumValueType[e.Name] = v.Name
			}
		case *ast.TypedefNode:
			c.definedTypedef[v.Name] = v
		case *ast.InterfaceNode:
			if len(v.Attributes) > 0 {
				c.definedInterface[v.Name] = v
			}
		case *ast.LibraryNode:
			defineIdent(c, v.Nodes)
		case *ast.CoClassNode:
			for _, ifc := range v.Interfaces {
				c.publicReturnedInterfaces[ifc.Name] = true
			}
		default:
		}
	}
}

func newDefContext(idls ...string) (*defContext, error) {
	c := &defContext{
		definedIdls:          make(map[string]interface{}),
		definedEnum:          make(map[string]*ast.EnumNode),
		definedEnumValueType: make(map[string]string),
		definedStruct:        make(map[string]*ast.StructNode),
		definedInterface:     make(map[string]*ast.InterfaceNode),
		definedTypedef:       make(map[string]*ast.TypedefNode),

		publicReturnedInterfaces: make(map[string]bool),
	}

	for _, fn := range idls {
		n := filepath.Base(fn)
		n = strings.TrimSuffix(n, ".idl")

		f, err := os.Open(fn)
		if err != nil {
			return nil, err
		}

		rootnodes := parser.Parse(f)
		defineIdent(c, rootnodes)
		c.definedIdls[n] = rootnodes
	}

	var queue []string

	for i := range c.publicReturnedInterfaces {
		queue = append(queue, i)
	}

	// TODO this is not temp way: clear the map to hide all com clients
	// since they are now moved to client com hub
	c.publicReturnedInterfaces = make(map[string]bool)

	for ; len(queue) > 0; queue = queue[1:] {
		ifc := c.definedInterface[queue[0]]

		if ifc == nil {
			continue
		}

		for _, m := range ifc.Methods {
			public := casee.IsPascalCase(goMethodName(m.Name))
			if public {
				for _, p := range m.Params {
					if isOutParam(p) {
						c.publicReturnedInterfaces[p.Type] = true
						queue = append(queue, p.Type)
					}
				}
			}
		}
	}

	// temp hack for typedef
	c.definedTypedef["FABRIC_STRING_PAIR"] = &ast.TypedefNode{
		Name: "FABRIC_STRING_PAIR",
		Type: "FABRIC_APPLICATION_PARAMETER",
	}

	c.definedTypedef["PCFABRIC_X509_NAME"] = &ast.TypedefNode{
		Name:        "PCFABRIC_X509_NAME",
		Type:        "FABRIC_X509_NAME",
		Indirection: 1,
	}

	c.definedTypedef["PCFABRIC_X509_ISSUER_NAME"] = &ast.TypedefNode{
		Name:        "PCFABRIC_X509_ISSUER_NAME",
		Type:        "FABRIC_X509_ISSUER_NAME",
		Indirection: 1,
	}

	c.definedTypedef["REFIID"] = &ast.TypedefNode{
		Name:        "REFIID",
		Type:        "GUID",
		Indirection: 1,
	}

	return c, nil
}
