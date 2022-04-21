package brute

import (
	"fmt"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"gopkg.in/mgo.v2"
	"time"
)

type mongodb struct {
	timeout time.Duration
}

func newBruteMongodb(L *lua.LState) service {
	opt := L.CheckTable(1)
	port := cast.ToInt(opt.RawGetString("port").String())

	sv := &mongodb{
		timeout: time.Duration(cast.ToInt(opt.RawGetString("timeout").String())),
	}
	if sv.timeout == 0 {
		logger.Errorf("mongodb timeout not set: %s , default 5", opt.RawGetString("timeout").String())
		sv.timeout = 5 * time.Second
	}
	return newService(L, sv, port)
}

func (m *mongodb) Name() string {
	return "mongodb"
}

func (m *mongodb) Login(ev *event) {
	url := fmt.Sprintf("mongodb://%v:%v@%v/%v", ev.user, ev.pass, ev.Server(), "admin")
	session, err := mgo.DialWithTimeout(url, m.timeout)

	if err == nil {
		defer session.Close()
		err = session.Ping()
		if err == nil {
			ev.stat = Succeed
			println("success", ev.user, ev.pass)
			//o.ev(ip, user, pass, port, "mongodb hit")
		} else {
			//println("xxx:",ev.user,ev.pass,err.Error())
			//logger.Error("ping mongodb err : %s", err.Error())
		}
	} else {
		//println(ev.user,ev.pass,err.Error())
		//logger.Errorf("connect mongodb err : %s",err.Error())
		//*userstop = true
		//*ipstopf = true
	}
}
