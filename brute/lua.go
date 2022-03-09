package brute

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xbase"
	"reflect"
)

var xEnv *xbase.EnvT

var typeof = reflect.TypeOf((*brute)(nil)).String()


/*
	local kfk = kafka.producer{}
	local b = crack.brute("crack").use("root" , "admin").pass("123456" , "780123").pipe(kfk)
	local b = crack.brute("crack").pass("123456" , "780123").pipe(kfk)

	//b.L("ip.txt")
	local b = crack.brute("crack").Use("user.txt").Pass("pass.txt").SOMBIE()
	b.host("www.baidu.com")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")
	//b.cidr("192.168.1.0/24" , "192.168.1.1-10" , "192.168.0.0/16")

	b.ssh(22,45865,10-100).skip(true)
	//b.ssh(tcp://127.0.0.1:22).ping(true).skip(true).proxy("sock5://root:123@192.168.1.1:8080")
	//b.mysql(3306).ping(true).skip(false)

	b.web("https://{host}/api/login?user={user}&pass={pass}}")
		.method("GET")
		.header{}
		.data("user={user}&pass={pass})
		.code(200 ,302 , 301)
		.filter("*ok|connection")
		.pipe(_(r)
			r.raw:regex("*login succeed")
			r.ok()
		end)

	b.sync()
	b.async()
	b.start()

	-----------------------------------------------------------
	{ip , io.reader , io.reader} , {ip , io.reader , io.reader}
    -----------------------------------------------------------
 */

func bruteL(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name, typeof)
	if proc.IsNil() {
		proc.Set(newBrute(cfg))

	} else {
		obj := proc.Data.(*brute)
		xEnv.Free(obj.cfg.co)
		obj.cfg = cfg
	}
	L.Push(proc)

	return 0
}

func Constructor(env *xbase.EnvT, ck lua.UserKV) {
	xEnv = env
	ck.Set("brute", lua.NewFunction(bruteL))
}
