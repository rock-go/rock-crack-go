package brute

import (
	"database/sql"
	"fmt"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"strings"
)

type mysql struct {
}

func newBruteMysql(L *lua.LState) service {
	opt := L.CheckTable(1)
	port := cast.ToInt(opt.RawGetString("port").String())

	sv := &mysql{
		//timeout: time.Duration(cast.ToInt(opt.RawGetString("timeout").String())),
	}
	return newService(L, sv, port)
}

func (m *mysql) Name() string {
	return "mysql"
}

func (m *mysql) Login(ev *event) {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", ev.user, ev.pass, ev.ip, ev.port, "mysql")
	db, err := sql.Open("mysql", dataSourceName)
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			println("success", ev.user, ev.pass)
			//o.ev(ip, user, pass, port, "mysql hit")

		} else {
			if !strings.Contains(err.Error(), "Error 1045") { //非密码错误问题
				logger.Errorf("mysql crack err , ip is baned : &s", err.Error())
				ev.stat = Fail
			}
		}
	}
	//ip有问题
}
