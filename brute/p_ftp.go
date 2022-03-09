package brute


import (
	 "github.com/jlaffaye/ftp"
	"github.com/rock-go/rock/lua"
	"strings"
	"time"
)

type Ftp struct {
	timeout  time.Duration
}

func newBruteFtp(L *lua.LState) service {
	return newService(L , &Ftp{timeout: time.Second})
}

func (f *Ftp) Name() string {
	return "ftp"
}

func (f *Ftp) Login(ev *event) {
	conn, err := ftp.DialTimeout(ev.Server() , f.timeout)

	if err != nil {
		ev.stat = Fail
		ev.banner = err.Error()
		return
	}

	err = conn.Login(ev.user, ev.pass)
	if err != nil {
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
}

func (ftp *Ftp) Index(L *lua.LState , key string) lua.LValue {
	return lua.LNil
}