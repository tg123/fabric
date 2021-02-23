package main

import (
	"fmt"

	"github.com/jd3nn1s/gomidl/ast"
	"github.com/pinzolo/casee"
)

func (g *generator) visitConstDef(n *ast.ConstdefNode) {

	t := g.toGolangType(n.Type, 0, false)
	v := n.Val

	if t == "string" {
		v = fmt.Sprintf(`"%v"`, v)
	}

	g.printfln("const %v = %v", casee.ToPascalCase(n.Name), v)
}
