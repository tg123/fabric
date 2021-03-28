package main

import (
	"fmt"
	"strings"
)

type callStubBuilder struct {
	uid          int
	callstub     map[string]int
	callbackstub map[string]int
	calls        map[int][]string
	callbacks    map[int][]string
}

func newCallStubBuilder() *callStubBuilder {
	return &callStubBuilder{
		callstub:     make(map[string]int),
		calls:        make(map[int][]string),
		callbackstub: make(map[string]int),
		callbacks:    make(map[int][]string),
	}
}

func (c *callStubBuilder) makeStub(stub map[string]int, funcs map[int][]string, paramTypes []string) int {
	k := strings.Join(paramTypes, ",")

	u, ok := stub[k]

	if !ok {
		c.uid++
		u = c.uid
		stub[k] = c.uid
		tmp := make([]string, len(paramTypes))
		copy(tmp, paramTypes)
		funcs[u] = tmp
	}

	return u
}

func (c *callStubBuilder) MakeCallStub(paramTypes []string) string {
	return c.CallStubName(c.makeStub(c.callstub, c.calls, paramTypes))
}

func (c *callStubBuilder) MakeCallbackStub(paramTypes []string) string {
	return c.CallbackStubName(c.makeStub(c.callbackstub, c.callbacks, paramTypes))
}

func (c *callStubBuilder) CallStubName(u int) string {
	return fmt.Sprintf("callStub%v", u)
}

func (c *callStubBuilder) CallbackStubName(u int) string {
	return fmt.Sprintf("createCallbackStub%v", u)
}
