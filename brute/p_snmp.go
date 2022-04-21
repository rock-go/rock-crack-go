package brute

import (
	"github.com/gosnmp/gosnmp"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"time"
)

type snmp struct {
	timeout time.Duration
}

func newBruteSnmp(L *lua.LState) service {
	opt := L.CheckTable(1)
	port := cast.ToInt(opt.RawGetString("port").String())

	sv := &snmp{
		timeout: time.Duration(cast.ToInt(opt.RawGetString("timeout").String())),
	}
	return newService(L, sv, port)
}

func (s *snmp) Name() string {
	return "snmp"
}

func (s *snmp) Login(ev *event) {
	gosnmp.Default.Target = ev.ip
	gosnmp.Default.Port = uint16(ev.port)
	gosnmp.Default.Community = ev.pass
	gosnmp.Default.Timeout = s.timeout

	err := gosnmp.Default.Connect()
	if err == nil {
		oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
		_, err := gosnmp.Default.Get(oids)
		if err == nil {
			//println(pass)
			//o.ev(ip, user, pass, port, "redis hit")
		} else {
			println(ev.pass, err.Error())
		}
	} else {
		println(ev.pass, err.Error())
	}
}
