package john

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xbase"
)

var xEnv *xbase.EnvT

/*
	local function handle(ev)
		ev.Put(true , true)
	end


	local john = crack.john{
		name = "shadow",
		dict = "share/dict/pass.dict",
		pipe = handle
	}

	local john = crack.john("shadow")
         .dict("share/dict/pass.dict")
         .pipe(handle)

	john.shadow("$1$xxx")
	john.md5("xxxx")
	john.sha256("xxx")
	john.sha512("xxxxxx")

 */

func newLuaCrackJohn(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name , fileTypeOf)
	if proc.IsNil() {
		proc.Set(newJohn(cfg))

	} else {
		obj := proc.Data.(*john)
		xEnv.Free(obj.cfg.co)
		obj.cfg = cfg
	}

	return 0
}

func Constructor(env *xbase.EnvT , ck lua.UserKV) {
	xEnv = env
	ck.Set("john" , lua.NewFunction(newLuaCrackJohn))
}
