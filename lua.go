package crack

import (
	"github.com/rock-go/rock-crack-go/john"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xbase"
)

func LuaInjectApi(env *xbase.EnvT) {
	ck := lua.NewUserKV()
	john.Constructor(env , ck)
	env.Global("crack" , ck)
}
