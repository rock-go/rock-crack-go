package john

import (
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
)

type config struct {
	speed int
	name  string
	dict  string

	co    *lua.LState
	pipe  *lua.LFunction
}

func (c config) verify() interface{} {
	return nil
}

func (c *config) Index(L *lua.LState , key string , val lua.LValue) {
	switch key {
	case "name":
		c.name = auxlib.CheckProcName(val , L)

	case "speed":
		c.speed = lua.IsInt(val)

	case "dict":
		c.dict = lua.IsString(val)

	default:
		L.RaiseError("invalid %s field", key)
		return
	}
}


func newConfig(L *lua.LState) *config {
	val := L.Get(1)
	cfg := &config{}

	switch val.Type() {
	case lua.LTString:
		cfg.name = auxlib.CheckProcName(val , L)

	case lua.LTTable:
		val.(*lua.LTable).Range(func(key string, val lua.LValue) { cfg.Index(L , key , val) })

	default:
		L.RaiseError("invalid config type must string or table , got %s" , val.Type().String())
		return nil
	}

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return cfg

}