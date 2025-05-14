package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("print", L.NewFunction(capturePrint))
	err := L.DoString(`local s = print("Hello from lua!")`)
	if err != nil {
		panic(err)
	}

	err = L.DoFile("cmd/tutorials/lua/test.lua")
	if err != nil {
		panic(err)
	}

}

func capturePrint(L *lua.LState) int {
	str := L.ToString(1)
	fmt.Println("Lua says:", str)
	return 0
}
