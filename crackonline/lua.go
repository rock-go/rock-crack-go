package crackonline

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xbase"
)

var xEnv *xbase.EnvT

/*
	local function handle(ev)
		ev.Put(true , true)
	end


	local online = crack.online{
		name   = "crack",
		iplist = "share/ip.txt",
        dict   = "share/pass.txt"
		pipe = handle
	}
    online.ssh(22)
    online.mysql(3306)

*/

func newLuaCrackOnline(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name, fileTypeOf)
	if proc.IsNil() {
		proc.Set(newJohn(cfg))

	} else {
		obj := proc.Data.(*online)
		xEnv.Free(obj.cfg.co)
		obj.cfg = cfg
	}
	L.Push(proc)
	return 1
}

func Constructor(env *xbase.EnvT, ck lua.UserKV) {
	xEnv = env
	ck.Set("online", lua.NewFunction(newLuaCrackOnline))
}
