package brute

import (
	"github.com/22ke/gordp/glog"
	gor "github.com/22ke/gordp/login"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"strings"
)

type rdp struct {
}

func newBruteRdp(L *lua.LState) service {
	opt := L.CheckTable(1)
	port := cast.ToInt(opt.RawGetString("port").String())

	sv := &sshd{
		//timeout: time.Duration(cast.ToInt(opt.RawGetString("timeout").String())),
	}
	return newService(L, sv, port)
}

func (r *rdp) Name() string {
	return "rdp"
}

func (r *rdp) Login(ev *event) {
	var err error
	g := gor.NewClient(ev.Server(), glog.NONE)
	//SSL协议登录测试
	err = g.LoginForSSL("", ev.user, ev.pass)
	if err == nil {
		return
	}
	//println(err.Error())
	if strings.Contains(err.Error(), "success") {
		//println("login success , ", target, " : ", user, " : ", pass)
		//o.ev(ip, user, pass, port, "rdp hit")
	}
}
