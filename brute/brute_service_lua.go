package brute

import "github.com/rock-go/rock/lua"

func (s *service) pingL(L *lua.LState) int {
	s.ping = L.IsTrue(1)
	return 0
}

func (s *service) skipL(L *lua.LState) int {
	s.skip = L.IsTrue(1)
	return 0
}

func (s *service) Index(L *lua.LState , key string) lua.LValue {
	if s.auth == nil {
		return lua.LNil
	}

	if lv := s.auth.Index(L , key); lv.Type() != lua.LTNil {
		return lv
	}

	switch key {

	case "ping":
		return lua.NewFunction(s.pingL)

	case "skip":
		return lua.NewFunction(s.skipL)
	}

	return lua.LNil
}