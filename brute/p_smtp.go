package brute

import (
	"crypto/tls"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	sm "net/smtp"
	"strings"
)

type smtp struct {
}

func newBruteSmtp(L *lua.LState) service {
	opt := L.CheckTable(1)
	port := cast.ToInt(opt.RawGetString("port").String())

	sv := &smtp{}
	return newService(L, sv, port)
}

func (s *smtp) Name() string {
	return "smtp"
}
func (s *smtp) Login(ev *event) {
	str := ev.Server()
	c, err := sm.Dial(str)
	if err != nil {
		//println("dial",err.Error())
		logger.Errorf("dial %s err : %s", ev.ip, err.Error())
		return
	}
	auth := sm.PlainAuth("", ev.user, ev.pass, ev.ip)

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: ev.ip, InsecureSkipVerify: true}
		if err = c.StartTLS(config); err != nil {
			//fmt.Println("call start tls")
			//return err
		}
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				if strings.Contains(err.Error(), "504 Unrecognized authentication type") {
					logger.Errorf("smtp crack error: 线程数量过多！ %s", err.Error())
				}
				//密码错误
			} else {
				//println(user,pass)
				//o.ev(ip, user, pass, port, "smtp hit")
			}
		}
	}
}
