package brute

import (
	"fmt"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type sshd struct {
	timeout time.Duration
}

func newBruteSsh(L *lua.LState) service {
	opt := L.CheckTable(1)
	port := cast.ToInt(opt.RawGetString("port").String())

	sv := &sshd{
		timeout: time.Duration(cast.ToInt(opt.RawGetString("timeout").String())),
	}
	return newService(L, sv, port)
}

func (s *sshd) Name() string {
	return "ssh"
}

func (s *sshd) Login(ev *event) {
	if s.timeout == 0 {
		s.timeout = time.Duration(3)
	}

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
	client, err := ssh.Dial("tcp", ev.Server(), cfg)
	fmt.Printf("ssh login user:%v,pass:%v\n", ev.user, ev.pass)
	if err == nil {
		println("连接成功：", ev.user, ":", ev.pass)
		defer client.Close()
		ev.stat = Succeed
		ev.banner = auxlib.B2S(client.ServerVersion())
	} else { //密码错误
		//println("密码错误: ",err.Error())
		ev.stat = Fail
		ev.banner = err.Error()
	}
}

func (s *sshd) Index(L *lua.LState, key string) lua.LValue {
	println("L:", L.GetTop())

	switch key {
	case "timeout":
		//s.timeout = time.Duration(3)
	}
	return lua.LNil
}
