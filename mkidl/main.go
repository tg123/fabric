package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pinzolo/casee"
)

func main() {
	fmt.Println(runtime.GOARCH)
	ctx, err := newDefContext(os.Args[1:]...)

	if err != nil {
		panic(err)
	}

	for idl, root := range ctx.definedIdls {

		g := newGenerator(ctx)
		g.walk(root)

		out, err := os.OpenFile(fmt.Sprintf("z_generated_%v.go", casee.ToFlatCase(idl)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

		if err != nil {
			panic(err)
		}

		g.dump(out)
	}

	for _, s := range []string{"dummy", "windows", "linux"} {

		g := newGenerator(ctx)
		out, err := os.OpenFile(fmt.Sprintf("z_generated_stub_%v.go", s), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

		if err != nil {
			panic(err)
		}

		g.visitStub(s)

		g.dump(out)
	}
}
