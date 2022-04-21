package brute

import (
	"github.com/jlaffaye/ftp"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"strings"
	"time"
)

type Ftp struct {
	timeout time.Duration
}

func newBruteFtp(L *lua.LState) service {
	val := L.CheckTable(1)
	port := cast.ToInt(val.RawGetString("port").String())

	e := &Ftp{
		timeout: time.Duration(cast.ToInt(val.RawGetString("timeout").String())),
	}
	if e.timeout == 0 {
		logger.Errorf("ftp timeout not set: %s , default 5", val.RawGetString("timeout").String())
		e.timeout = 5 * time.Second
	}

	println("timeout: ", e.timeout)
	return newService(L, e, port)
}

func (f *Ftp) Name() string {
	return "ftp"
}

func (f *Ftp) Login(ev *event) {
	conn, err := ftp.DialTimeout(ev.Server(), f.timeout)

	if err != nil {
		ev.stat = Fail
		ev.banner = err.Error()
		return
	}

	err = conn.Login(ev.user, ev.pass)
	if err != nil {
		//println("fail \n", ev.user,ev.pass)
		banner := err.Error()
		if strings.Contains(ev.banner, "Permission denied") {
			ev.stat = Denied
		} else {
			ev.stat = Fail
		}

		ev.banner = banner
		return
	}
	defer conn.Logout()

	ev.stat = Succeed
	println("success \n", ev.user, ev.pass)
}
