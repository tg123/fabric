package main

import (
	"fmt"
	"hash/maphash"
	"strings"
)

type callStubBuilder struct {
	// uid          int

	callstub     map[string]uint64
	callbackstub map[string]uint64
	calls        map[uint64][]string
	callbacks    map[uint64][]string
}

func newCallStubBuilder() *callStubBuilder {
	return &callStubBuilder{
		callstub:     make(map[string]uint64),
		calls:        make(map[uint64][]string),
		callbackstub: make(map[string]uint64),
		callbacks:    make(map[uint64][]string),
	}
}

func (c *callStubBuilder) makeStub(stub map[string]uint64, funcs map[uint64][]string, paramTypes []string) uint64 {
	k := strings.Join(paramTypes, ",")

	u, ok := stub[k]

	if !ok {
		// c.uid++
		var h maphash.Hash
		h.WriteString(k)
		u = h.Sum64()

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

func (c *callStubBuilder) CallStubName(u uint64) string {
	return fmt.Sprintf("callStub%v", u)
}

func (c *callStubBuilder) CallbackStubName(u uint64) string {
	return fmt.Sprintf("createCallbackStub%v", u)
}
