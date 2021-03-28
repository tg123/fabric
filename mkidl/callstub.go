package main

import (
	"fmt"
	"hash/adler32"
	"strings"
)

type callStubBuilder struct {
	// uid          int

	callstub     map[string]uint32
	callbackstub map[string]uint32
	calls        map[uint32][]string
	callbacks    map[uint32][]string
}

func newCallStubBuilder() *callStubBuilder {
	return &callStubBuilder{
		callstub:     make(map[string]uint32),
		calls:        make(map[uint32][]string),
		callbackstub: make(map[string]uint32),
		callbacks:    make(map[uint32][]string),
	}
}

func (c *callStubBuilder) makeStub(stub map[string]uint32, funcs map[uint32][]string, paramTypes []string) uint32 {
	k := strings.Join(paramTypes, ",")

	u, ok := stub[k]

	if !ok {
		// c.uid++
		u = adler32.Checksum([]byte(k))

		for {
			_, ok := funcs[u]
			if !ok {
				break
			}

			u++
		}

		stub[k] = u
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

func (c *callStubBuilder) CallStubName(u uint32) string {
	return fmt.Sprintf("callStub%v", u)
}

func (c *callStubBuilder) CallbackStubName(u uint32) string {
	return fmt.Sprintf("createCallbackStub%v", u)
}
