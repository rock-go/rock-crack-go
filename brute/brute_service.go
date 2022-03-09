package brute

import (
	"github.com/rock-go/rock/lua"
	"net"
)

type Authenticate interface {
	Name() string
	Login(*event)
	Index(*lua.LState , string) lua.LValue
}

type service struct {
	skip    bool
	port    int
	ping    bool
	auth    Authenticate
}

func newService(L *lua.LState , auth Authenticate) service {
	return service{
		port: toPort(L) ,
		skip: true,
		ping: true,
		auth: auth,
	}
}

func (s *service) Ping(ip net.IP) bool {
	if !s.ping {
		return true
	}

	//写ping逻辑
	return true
}

func (s *service) Do(ev *event) {
	s.auth.Login(ev)
}