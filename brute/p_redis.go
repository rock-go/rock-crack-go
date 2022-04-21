package brute

import (
	red "github.com/go-redis/redis"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"time"
)

type redis struct {
	timeout time.Duration
}

func newBruteRedis(L *lua.LState) service {
	opt := L.CheckTable(1)
	port := cast.ToInt(opt.RawGetString("port").String())

	sv := &redis{
		timeout: time.Duration(cast.ToInt(opt.RawGetString("timeout").String())),
	}
	return newService(L, sv, port)
}

func (r *redis) Name() string {
	return "redis"
}

func (r *redis) Login(ev *event) {
	opt := red.Options{Addr: ev.Server(),
		Password: ev.pass, DB: 0, DialTimeout: r.timeout}
	client := red.NewClient(&opt)
	defer client.Close()
	_, err := client.Ping().Result()
	if err == nil {
		//println(pass)
		//o.ev(ip, user, pass, port, "redis hit")

	} else {
		//println(pass,err.Error())
	}
}
