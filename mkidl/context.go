package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/jd3nn1s/gomidl/parser"
)

type defContext struct {
	definedIdls map[string]interface{}

	// for forward lookup
	definedEnum          map[string]*ast.EnumNode
	definedEnumValueType map[string]string
	definedStruct        map[string]*ast.StructNode
	definedInterface     map[string]*ast.InterfaceNode
	definedTypedef       map[string]*ast.TypedefNode
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
		default:
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
}

func newDefContext(idls ...string) (*defContext, error) {
	c := &defContext{
		definedIdls:          make(map[string]interface{}),
		definedEnum:          make(map[string]*ast.EnumNode),
		definedEnumValueType: make(map[string]string),
		definedStruct:        make(map[string]*ast.StructNode),
		definedInterface:     make(map[string]*ast.InterfaceNode),
		definedTypedef:       make(map[string]*ast.TypedefNode),
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

	return c, nil
}
