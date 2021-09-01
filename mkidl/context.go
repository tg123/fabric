package main

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/jd3nn1s/gomidl/parser"
)

type defContext struct {
	definedIdls map[string]interface{}

	// for forward lookup
	definedEnum           map[string]*ast.EnumNode
	definedEnumValueType  map[string]string
	definedStruct         map[string]*ast.StructNode
	definedStructEx       map[string][]string
	definedStructExParent map[string]string
	definedInterface      map[string]*ast.InterfaceNode
	definedTypedef        map[string]*ast.TypedefNode

	stubBuilder *callStubBuilder
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
		default:
		}
	}
}

func newDefContext(idls ...string) (*defContext, error) {
	c := &defContext{
		definedIdls:           make(map[string]interface{}),
		definedEnum:           make(map[string]*ast.EnumNode),
		definedEnumValueType:  make(map[string]string),
		definedStruct:         make(map[string]*ast.StructNode),
		definedStructEx:       make(map[string][]string),
		definedStructExParent: make(map[string]string),
		definedInterface:      make(map[string]*ast.InterfaceNode),
		definedTypedef:        make(map[string]*ast.TypedefNode),

		stubBuilder: newCallStubBuilder(),
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

	c.definedStruct["TEST_COMMAND_QUERY_RESULT_LIST"].Fields[1].Type = "TEST_COMMAND_QUERY_RESULT_ITEM"

	for s := range c.definedStruct {
		p := strings.Split(s, "_EX")
		if len(p) == 2 {
			if _, err := strconv.Atoi(p[1]); err != nil {
				continue
			}
			parent := p[0]
			c.definedStructExParent[s] = parent
			c.definedStructEx[parent] = append(c.definedStructEx[parent], s)
		}
	}

	// what sort if doing an fucking interview
	for _, ex := range c.definedStructEx {
		sort.Slice(ex, func(i, j int) bool {
			return strings.Compare(ex[i], ex[j]) < 0
		})
	}

	return c, nil
}
