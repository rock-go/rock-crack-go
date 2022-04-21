package brute

import (
	"fmt"
	"github.com/rock-go/rock/lua"
)

func (s *service) pingL(L *lua.LState) int {
	s.ping = L.IsTrue(1)
	return 1
}

func (s *service) skipL(L *lua.LState) int {
	s.skip = L.IsTrue(1)
	return 1
}

func (s *service) Index(L *lua.LState, key string) lua.LValue {
	if s.auth == nil {
		return lua.LNil
	}
	str, _ := L.Get(1).AssertString()
	fmt.Printf("%v\n", str)
	fmt.Printf("%v\n", L.Get(2))

	println("LL", L.GetTop())
	//if lv := s.auth.Index(L , key); lv.Type() != lua.LTNil {
	//	return lv
	//}

	switch key {

	case "ping":
		return lua.NewFunction(s.pingL)

	case "skip":
		return lua.NewFunction(s.skipL)
	}

	//s.auth.Index(L,key)

	//L.Push(L.NewLightUserData(s))
	return lua.LNil

	//return lua.LNil
}
