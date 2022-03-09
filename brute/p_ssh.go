package brute

import (
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type sshd struct {
	timeout  time.Duration
}

func newBruteSsh(L *lua.LState) service {
	sv := &sshd{}
	return newService(L , sv)
}

func (s *sshd) Name() string {
	return "ssh"
}

func (s *sshd) Login(ev *event) {
	cfg := &ssh.ClientConfig{
		User: ev.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(ev.pass),
		},
		Timeout: s.timeout * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", ev.Server() , cfg)
	if err == nil {
		defer client.Close()
		ev.stat = Succeed
		ev.banner = auxlib.B2S(client.ServerVersion())
	} else { //密码错误
		ev.stat = Fail
		ev.banner = err.Error()
	}
}

func (s *sshd) Index(L *lua.LState , key string) lua.LValue {
	return nil
}